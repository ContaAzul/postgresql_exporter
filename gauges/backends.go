package gauges

import "github.com/prometheus/client_golang/prometheus"

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

func (g *Gauges) BackendsStatus() *prometheus.GaugeVec {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_backends",
		Help:        "Active database connections",
		ConstLabels: g.labels,
	}
	var gauge = prometheus.NewGaugeVec(opts, []string{"status"})
	for _, status := range []string{"active", "idle", "idle in transaction"} {
		g.from(
			gauge.With(prometheus.Labels{"status": status}),
			`
				SELECT COUNT(*)
				FROM pg_stat_activity
				WHERE datname = current_database()
				AND state = $1
				AND pid <> pg_backend_pid()
			`,
			status,
		)
	}
	return gauge
}

func (g *Gauges) WaitingBackends() prometheus.Gauge {
	var query = `
		SELECT COUNT(*)
		FROM pg_stat_activity
		WHERE datname = current_database()
		AND waiting is true
	`
	if isPG96(g.version()) {
		query = `
			SELECT COUNT(*)
			FROM pg_stat_activity
			WHERE datname = current_database()
			AND wait_event is not null
		`
	}
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_waiting_backends",
			Help:        "Database connections waiting on a Lock",
			ConstLabels: g.labels,
		},
		query,
	)
}
