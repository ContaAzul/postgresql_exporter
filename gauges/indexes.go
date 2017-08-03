package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) UnusedIndexes() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_unused_indexes",
			Help:        "Dabatase unused indexes count",
			ConstLabels: g.labels,
		},
		`
			SELECT COUNT(*)
			FROM pg_stat_user_indexes ui
			JOIN pg_index i ON ui.indexrelid = i.indexrelid
			WHERE NOT i.indisunique
			AND ui.idx_scan < 100
		`,
	)
}
