package gauges

import (
	"fmt"

	"github.com/ContaAzul/postgresql_exporter/postgres"
	"github.com/prometheus/client_golang/prometheus"
)

// ReplicationStatus returns a prometheus gauge for the PostgreSQL
// replication status
func (g *Gauges) ReplicationStatus() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_status",
			Help:        "Returns 1 if in recovery and replay is paused, 0 if OK and -1 if not in recovery",
			ConstLabels: g.labels,
		},
		fmt.Sprintf(`
			SELECT
			CASE
				WHEN pg_is_in_recovery() is true
				THEN
					CASE
						WHEN %s() is true
						THEN 1
						ELSE 0
					END
			ELSE
				-1
			END`,
			postgres.Version(g.version()).IsWalReplayPausedFunctionName(),
		),
	)
}

// StreamingWALs returns a prometheus gauge for the count of WALs
// in streaming state
func (g *Gauges) StreamingWALs() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_streaming_wals",
			Help:        "Returns the count of WALs in streaming state",
			ConstLabels: g.labels,
		},
		`
			SELECT count(state)
			FROM pg_stat_replication
			WHERE state='streaming'
		`,
	)
}

// ReplicationDelayInSeconds returns a prometheus gauge for the database replication
// lag in seconds
func (g *Gauges) ReplicationDelayInSeconds() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_delay_seconds",
			Help:        "Dabatase replication delay in seconds",
			ConstLabels: g.labels,
		},
		"SELECT coalesce(extract(epoch from now() - pg_last_xact_replay_timestamp()), 0) AS replication_delay",
	)
}

// ReplicationDelayInBytes returns a prometheus gauge for the database replication
// lag in bytes
func (g *Gauges) ReplicationDelayInBytes() prometheus.Gauge {
	version := postgres.Version(g.version())

	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_delay_bytes",
			Help:        "Dabatase replication delay in bytes",
			ConstLabels: g.labels,
		},
		fmt.Sprintf(`
			SELECT COALESCE(ABS(%s(%s(), %s())), 0) AS replication_delay_bytes`,
			version.WalLsnDiffFunctionName(),
			version.LastWalReceivedLsnFunctionName(),
			version.LastWalReplayedLsnFunctionName(),
		),
	)
}
