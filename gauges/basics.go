package gauges

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

func Up(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	return newGauge(
		db,
		prometheus.GaugeOpts{
			Name:        "postgresql_up",
			Help:        "Dabatase is up and accepting connections",
			ConstLabels: labels,
		},
		"SELECT 1",
	)
}

func Size(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	return newConvertedGauge(
		db,
		prometheus.GaugeOpts{
			Name:        "postgresql_db_size_mb",
			Help:        "Dabatase size in mbs",
			ConstLabels: labels,
		},
		"SELECT pg_database_size(current_database())",
		func(result float64) float64 {
			return result / (1024 * 1024)
		},
	)
}

func Deadlocks(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	return newGauge(
		db,
		prometheus.GaugeOpts{
			Name:        "postgresql_deadlock",
			Help:        "Number of deadlocks in the last 2m",
			ConstLabels: labels,
		},
		`
			SELECT count(*) FROM pg_locks bl
			JOIN pg_stat_activity a
			ON a.pid = bl.pid JOIN pg_locks kl
			ON kl.transactionid = bl.transactionid
			AND kl.pid != bl.pid JOIN pg_stat_activity ka
			ON ka.pid = kl.pid WHERE NOT bl.granted
			AND (ka.query_start >= (now() - interval '2 minutes')
			OR a.query_start >= (now() - interval '2 minutes'))
			AND a.datname = current_database()
		`,
	)
}
