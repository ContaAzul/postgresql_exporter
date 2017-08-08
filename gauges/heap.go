package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) HeapBlocksRead() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_heap_blks_read_sum",
			Help:        "Sum of the number of disk blocks read from all tables",
			ConstLabels: g.labels,
		},
		"SELECT coalesce(sum(heap_blks_read), 0) FROM pg_statio_user_tables",
	)
}

func (g *Gauges) HeapBlocksHit() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_heap_blks_hit_sum",
			Help:        "Sum of the number of buffer hits on all tables",
			ConstLabels: g.labels,
		},
		"SELECT coalesce(sum(heap_blks_hit), 0) FROM pg_statio_user_tables",
	)
}
