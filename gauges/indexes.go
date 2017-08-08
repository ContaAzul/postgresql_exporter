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

func (g *Gauges) IndexBlocksRead() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_index_blks_read_sum",
			Help:        "Sum of the number of disk blocks read from all public indexes",
			ConstLabels: g.labels,
		},
		`
			SELECT coalesce(sum(idx_blks_read), 0)
			FROM pg_statio_all_indexes
			WHERE schemaname = 'public'
		`,
	)
}

func (g *Gauges) IndexBlocksHit() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_index_blks_hit_sum",
			Help:        "Sum of the number of buffer hits on all public indexes",
			ConstLabels: g.labels,
		},
		`
			SELECT coalesce(sum(idx_blks_hit), 0)
			FROM pg_statio_all_indexes
			WHERE schemaname = 'public'
		`,
	)
}
