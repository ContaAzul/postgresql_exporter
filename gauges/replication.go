package gauges

import "github.com/prometheus/client_golang/prometheus"

func (g *Gauges) ReplicationStatus() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_status",
			Help:        "Returns 1 if in recovery and replay is paused, 0 if OK and -1 if not in recovery",
			ConstLabels: g.labels,
		},
		`
			SELECT
			CASE
				WHEN pg_is_in_recovery() is true
				THEN
					CASE
						WHEN pg_is_xlog_replay_paused() is true
						THEN 1
						ELSE 0
					END
			ELSE
				-1
			END
		`,
	)
}

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

func (g *Gauges) ReplicationLag() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_lag",
			Help:        "Dabatase replication lag",
			ConstLabels: g.labels,
		},
		`
			SELECT COALESCE(
				CASE
					WHEN pg_is_in_recovery() is true
					THEN
					CASE
						WHEN pg_last_xlog_receive_location() = pg_last_xlog_replay_location()
						THEN 0
					ELSE
						EXTRACT (EPOCH FROM now() - pg_last_xact_replay_timestamp())
					END
				END
			, 0)
		`,
	)
}
