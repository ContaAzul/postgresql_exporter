package gauges

import (
	"time"

	"github.com/ContaAzul/postgresql_exporter/postgres"
	"github.com/prometheus/client_golang/prometheus"
)

// ConnectedBackends returns the number of backends currently connected to database
func (g *Gauges) ConnectedBackends() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_total",
			Help:        "Number of backends currently connected to database",
			ConstLabels: g.labels,
		},
		"SELECT numbackends FROM pg_stat_database WHERE datname = current_database()",
	)
}

// MaxBackends returns the maximum number of concurrent connections in the database
func (g *Gauges) MaxBackends() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_max_backends",
			Help:        "Maximum number of concurrent connections in the database",
			ConstLabels: g.labels,
		},
		"SELECT setting::numeric FROM pg_settings WHERE name = 'max_connections'",
	)
}

// CurrentBackends returns the number of backends currently connected to all databases 
func (g *Gauges) CurrentBackends() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_current_backends",
			Help:        "Current number of concurrent connections in all databases",
			ConstLabels: g.labels,
		},
		"SELECT sum(numbackends) FROM pg_stat_database;",
	)
}

type backendsByState struct {
	Total float64 `db:"total"`
	State string  `db:"state"`
}

// BackendsByState returns the number of backends currently connected to database by state
func (g *Gauges) BackendsByState() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_by_state_total",
			Help:        "Number of backends currently connected to database by state",
			ConstLabels: g.labels,
		},
		[]string{"state"},
	)

	const backendsByStateQuery = `
		SELECT COUNT(*) AS total, state
		FROM pg_stat_activity
		WHERE datname = current_database()
		GROUP BY state
	`

	go func() {
		for {
			gauge.Reset()
			var backendsByState []backendsByState
			if err := g.query(backendsByStateQuery, &backendsByState, emptyParams); err == nil {
				for _, row := range backendsByState {
					gauge.With(prometheus.Labels{
						"state": row.State,
					}).Set(row.Total)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

type backendsByUserAndClientAddress struct {
	Total      float64 `db:"total"`
	User       string  `db:"usename"`
	ClientAddr string  `db:"client_addr"`
}

// BackendsByUserAndClientAddress returns the number of backends currently connected
// to database by user and client address
func (g *Gauges) BackendsByUserAndClientAddress() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_by_user_total",
			Help:        "Number of backends currently connected to database by user and client address",
			ConstLabels: g.labels,
		},
		[]string{"user", "client_addr"},
	)

	const backendsByUserAndClientAddressQuery = `
		SELECT
		  COUNT(*) AS total,
		  usename,
		  COALESCE(client_addr, '::1') AS client_addr
		FROM pg_stat_activity
		WHERE datname = current_database()
		GROUP BY usename, client_addr
	`

	go func() {
		for {
			gauge.Reset()
			var backendsByUserAndClientAddress []backendsByUserAndClientAddress
			if err := g.query(backendsByUserAndClientAddressQuery, &backendsByUserAndClientAddress, emptyParams); err == nil {
				for _, row := range backendsByUserAndClientAddress {
					gauge.With(prometheus.Labels{
						"user":        row.User,
						"client_addr": row.ClientAddr,
					}).Set(row.Total)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

type backendsByWaitEventType struct {
	Total         float64 `db:"total"`
	WaitEventType string  `db:"wait_event_type"`
}

func (g *Gauges) backendsByWaitEventTypeQuery() string {
	if postgres.Version(g.version()).IsEqualOrGreaterThan96() {
		return `
			SELECT
			  COUNT(*) AS total,
			  wait_event_type
			FROM pg_stat_activity
			WHERE wait_event_type IS NOT NULL
			  AND datname = current_database()
			GROUP BY wait_event_type
		`
	}
	return `
		SELECT
		  COUNT(*) as total,
		  'Lock' as wait_event_type
		FROM pg_stat_activity
		WHERE datname = current_database()
		  AND waiting is true
	`
}

// BackendsByWaitEventType returns the number of backends currently waiting on some event
func (g *Gauges) BackendsByWaitEventType() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_by_wait_event_type_total",
			Help:        "Number of backends currently waiting on some event",
			ConstLabels: g.labels,
		},
		[]string{"wait_event_type"},
	)

	go func() {
		for {
			gauge.Reset()
			var backendsByWaitEventType []backendsByWaitEventType
			if err := g.query(g.backendsByWaitEventTypeQuery(),
				&backendsByWaitEventType, emptyParams); err == nil {
				for _, row := range backendsByWaitEventType {
					gauge.With(prometheus.Labels{
						"wait_event_type": row.WaitEventType,
					}).Set(row.Total)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}
