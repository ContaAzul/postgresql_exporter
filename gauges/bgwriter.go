package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) RequestedCheckpoints() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_checkpoints_requested",
			Help:        "Number of requested checkpoints that have been performed",
			ConstLabels: g.labels,
		},
		`
			SELECT checkpoints_req
			FROM pg_stat_bgwriter
		`,
	)
}

func (g *Gauges) ScheduledCheckpoints() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_checkpoints_scheduled",
			Help:        "Number of scheduled checkpoints that have been performed",
			ConstLabels: g.labels,
		},
		`
			SELECT checkpoints_timed
			FROM pg_stat_bgwriter
		`,
	)
}

func (g *Gauges) BufferOversize() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_checkpoints_buffer_oversize",
			Help:        "Number of times the background writer stopped a cleaning scan because it had written too many buffers",
			ConstLabels: g.labels,
		},
		`
			SELECT maxwritten_clean
			FROM pg_stat_bgwriter
		`,
	)
}

func (g *Gauges) BuffersWritten() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_checkpoints_buffer_written",
			Help:        "Number of buffers written directly by a backend",
			ConstLabels: g.labels,
		},
		`
			SELECT buffers_backend
			FROM pg_stat_bgwriter
		`,
	)
}
