package gauges

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

func IdleSessions(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	return newGauge(
		db,
		prometheus.GaugeOpts{
			Name:        "postgresql_db_idle_sessions",
			Help:        "Dabatase idle sessions",
			ConstLabels: labels,
		},
		`
			SELECT count(*)
			FROM pg_stat_activity
			WHERE datname = current_database()
			AND pid <> pg_backend_pid()
			AND state_change < current_timestamp - INTERVAL '5' MINUTE
		`,
	)
}
