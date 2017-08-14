package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) Tps() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_tps",
			Help:        "Transactions per second from current database",
			ConstLabels: g.labels,
		},
		`
			SELECT xact_commit + xact_rollback as tps
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}
