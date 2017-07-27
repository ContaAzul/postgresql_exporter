package gauges

import (
	"database/sql"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

func ReplicationLag(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_db_replication_lag",
		Help:        "Dabatase replication lag",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			var result float64
			if err := db.QueryRow(`
				SELECT COALESCE(
					CASE
						WHEN pg_last_xlog_receive_location() = pg_last_xlog_replay_location()
						THEN 0
					ELSE
						EXTRACT (EPOCH FROM now() - pg_last_xact_replay_timestamp())
					END
				, 0)
			`).Scan(&result); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
			}
			return result
		},
	)
}
