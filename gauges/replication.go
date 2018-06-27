package gauges

import (
	"fmt"

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
			g.defineFunction("pg_is_wal_replay_paused", "pg_is_xlog_replay_paused"),
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

// ReplicationLag returns a prometheus gauge for the database replication
// lag in milliseconds
func (g *Gauges) ReplicationLag() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_lag",
			Help:        "Dabatase replication lag",
			ConstLabels: g.labels,
		},
		fmt.Sprintf(`
			SELECT COALESCE(
				CASE
					WHEN pg_is_in_recovery() is true
					THEN
					CASE
						WHEN %s() = %s()
						THEN 0
					ELSE
						EXTRACT (EPOCH FROM now() - pg_last_xact_replay_timestamp())
					END
				END
			, 0)`,
			g.defineFunction("pg_last_wal_receive_lsn", "pg_last_xlog_receive_location"),
			g.defineFunction("pg_last_wal_replay_lsn", "pg_last_xlog_replay_location"),
		),
	)
}

func (g *Gauges) defineFunction(pg10Function, pg9xFunction string) string {
	if isPG10(g.version()) {
		return pg10Function
	}
	return pg9xFunction
}
