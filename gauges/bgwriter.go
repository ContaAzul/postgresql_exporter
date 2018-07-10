package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

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

type buffersWritten struct {
	checkpoint float64 `db:"buffers_checkpoint"`
	bgWriter   float64 `db:"buffers_clean"`
	backend    float64 `db:"buffers_backend"`
}

var buffersWrittenQuery = `
	SELECT 
	  buffers_checkpoint,
	  buffers_clean,
	  buffers_backend
	FROM pg_stat_bgwriter
`

// BuffersWritten returns the number of buffers written directly by a backend,
// by the background writer and during checkpoints
func (g *Gauges) BuffersWritten() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_buffers_written",
			Help:        "table scans statistics",
			ConstLabels: g.labels,
		},
		[]string{"written_by"},
	)

	go func() {
		for {
			var buffersWritten []buffersWritten
			if err := g.query(buffersWrittenQuery, &buffersWritten, emptyParams); err == nil {
				for _, writtenBy := range buffersWritten {
					gauge.With(prometheus.Labels{
						"written_by": "checkpoint",
					}).Set(writtenBy.checkpoint)
					gauge.With(prometheus.Labels{
						"written_by": "bgwriter",
					}).Set(writtenBy.bgWriter)
					gauge.With(prometheus.Labels{
						"written_by": "backend",
					}).Set(writtenBy.backend)
				}
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}
