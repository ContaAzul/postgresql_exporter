package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) TransactionsSum() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_transactions_sum",
			Help:        "Sum of all transactions in the database",
			ConstLabels: g.labels,
		},
		`
			SELECT xact_commit + xact_rollback
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}
