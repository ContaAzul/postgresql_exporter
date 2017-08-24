package gauges

import (
	"time"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

func (g *Gauges) Backends() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_total",
			Help:        "Total database backends",
			ConstLabels: g.labels,
		},
		`
			SELECT numbackends
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}

func (g *Gauges) MaxBackends() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_max_backends",
			Help:        "Maximum database backends (per postmaster)",
			ConstLabels: g.labels,
		},
		`
			SELECT setting::numeric
			FROM pg_settings
			WHERE name = 'max_connections'
		`,
	)
}

type backendStatus struct {
	Count      float64 `db:"count"`
	User       string  `db:"usename"`
	ClientAddr string  `db:"client_addr"`
	State      string  `db:"state"`
}

func (g *Gauges) BackendsStatus() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_backends_count",
			Help:        "Count of connections by state",
			ConstLabels: g.labels,
		},
		[]string{"status", "user", "client_addr"},
	)
	if !g.isSuperuser {
		log.Warn("postgresql_backends_count disabled because it requires a superuser to see queries from other users")
		return gauge
	}
	const backendsQuery = `
		SELECT COUNT(*) AS count, state, usename, client_addr
		FROM pg_stat_activity
		WHERE datname = current_database()
		AND client_addr IS NOT NULL
		GROUP BY state, usename, client_addr
	`
	go func() {
		for _, q := range []string{
			backendsQuery,
			g.waitingBackendsQuery(),
		} {
			var statuteses []backendStatus
			g.query(q, &statuteses, emptyParams)
			for _, status := range statuteses {
				gauge.With(prometheus.Labels{
					"status":      status.State,
					"user":        status.User,
					"client_addr": status.ClientAddr,
				}).Set(status.Count)
			}
		}
		time.Sleep(g.interval)
	}()
	return gauge
}

func (g *Gauges) waitingBackendsQuery() string {
	if isPG96(g.version()) {
		return `
			SELECT COUNT(*) as count, 'waiting' as state, usename
			FROM pg_stat_activity
			WHERE datname = current_database()
			AND wait_event is not null
			GROUP BY usename
		`
	}
	return `
		SELECT COUNT(*) as count, 'waiting' as state, usename
		FROM pg_stat_activity
		WHERE datname = current_database()
		AND waiting is true
		GROUP BY usename
	`
}
