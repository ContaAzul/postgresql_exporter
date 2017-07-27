package gauges

import (
	"database/sql"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

func IdleSessions(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_db_idle_sessions",
		Help:        "Dabatase idle sessions",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			var result float64
			if err := db.QueryRow(`
				SELECT count(*)
				FROM pg_stat_activity
				WHERE datname = current_database()
				AND pid <> pg_backend_pid()
				AND state_change < current_timestamp - INTERVAL '5' MINUTE
			`).Scan(&result); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
			}
			return result
		},
	)
}
