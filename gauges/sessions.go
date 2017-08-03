package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) IdleSessions() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_idle_sessions",
			Help:        "Dabatase idle sessions",
			ConstLabels: g.labels,
		},
		`
			SELECT count(*)
			FROM pg_stat_activity
			WHERE datname = current_database()
			AND pid <> pg_backend_pid()
			AND state_change < current_timestamp - INTERVAL '5' MINUTE
		`,
	)
}
