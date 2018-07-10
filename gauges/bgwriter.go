package gauges

import "github.com/prometheus/client_golang/prometheus"

// RequestedCheckpoints returns the number of requested checkpoints that have been performed
func (g *Gauges) RequestedCheckpoints() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_checkpoints_requested",
			Help:        "Number of requested checkpoints that have been performed",
			ConstLabels: g.labels,
		},
		"SELECT checkpoints_req FROM pg_stat_bgwriter",
	)
}

// ScheduledCheckpoints returns the number of scheduled checkpoints that have been performed
func (g *Gauges) ScheduledCheckpoints() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_checkpoints_scheduled",
			Help:        "Number of scheduled checkpoints that have been performed",
			ConstLabels: g.labels,
		},
		"SELECT checkpoints_timed FROM pg_stat_bgwriter",
	)
}

// BuffersMaxWrittenClean returns the number of times the background writer stopped a
// cleaning scan because it had written too many buffers
func (g *Gauges) BuffersMaxWrittenClean() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_buffers_maxwritten_clean",
			Help:        "Number of times the background writer stopped a cleaning scan because it had written too many buffers",
			ConstLabels: g.labels,
		},
		"SELECT maxwritten_clean FROM pg_stat_bgwriter",
	)
}

// BuffersWrittenByCheckpoints returns the number of buffers written during checkpoints
func (g *Gauges) BuffersWrittenByCheckpoints() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_buffers_checkpoint",
			Help:        "Number of buffers written during checkpoints",
			ConstLabels: g.labels,
		},
		"SELECT buffers_checkpoint FROM pg_stat_bgwriter",
	)
}

// BuffersWrittenByBgWriter returns the number of buffers written by the background writer
func (g *Gauges) BuffersWrittenByBgWriter() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_buffers_clean",
			Help:        "Number of buffers written by the background writer",
			ConstLabels: g.labels,
		},
		"SELECT buffers_clean FROM pg_stat_bgwriter",
	)
}

// BuffersWrittenByBackend returns the number of buffers written directly by a backend
func (g *Gauges) BuffersWrittenByBackend() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_buffers_backend",
			Help:        "Number of buffers written directly by a backend",
			ConstLabels: g.labels,
		},
		"SELECT buffers_backend FROM pg_stat_bgwriter",
	)
}
