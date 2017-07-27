package gauges

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

func ReplicationLag(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	return newGauge(
		db,
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_lag",
			Help:        "Dabatase replication lag",
			ConstLabels: labels,
		},
		`
			SELECT COALESCE(
				CASE
					WHEN pg_last_xlog_receive_location() = pg_last_xlog_replay_location()
					THEN 0
				ELSE
					EXTRACT (EPOCH FROM now() - pg_last_xact_replay_timestamp())
				END
			, 0)
		`,
	)
}
