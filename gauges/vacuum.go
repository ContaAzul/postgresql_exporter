package gauges

import (
	"github.com/prometheus/client_golang/prometheus"
)

// UnvacuumedTransactions returns the number of unvacuumed transactions
func (g *Gauges) UnvacuumedTransactions() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_unvacuumed_transactions_total",
			Help:        "Number of unvacuumed transactions",
			ConstLabels: g.labels,
		},
		"SELECT age(datfrozenxid) FROM pg_database WHERE datname = current_database()",
	)
}
