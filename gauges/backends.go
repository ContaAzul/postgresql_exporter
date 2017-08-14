package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (g *Gauges) Backends() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_totalbackends",
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
	Count float64 `db:"count"`
	User  string  `db:"usename"`
	State string  `db:"state"`
}

func (g *Gauges) BackendsStatus() *prometheus.GaugeVec {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_backends",
		Help:        "Active database connections",
		ConstLabels: g.labels,
	}
	var gauge = prometheus.NewGaugeVec(opts, []string{"status", "user"})
	const backendsQuery = `
		SELECT COUNT(*) as count, state, usename
		FROM pg_stat_activity
		WHERE datname = current_database()
		GROUP BY state, usename
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
					"status": status.State,
					"user":   status.User,
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
