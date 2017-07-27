package gauges

import (
	"database/sql"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

func Up(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_up",
		Help:        "Dabatase is up and accepting connections",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			if _, err := db.Query(`SELECT 1`); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
				return 1
			}
			return 0
		},
	)
}

func Size(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_db_size_mb",
		Help:        "Dabatase size in mbs",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			var result float64
			if err := db.QueryRow(`
				SELECT pg_database_size(current_database())
			`).Scan(&result); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
			}
			return result / (1024 * 1024)
		},
	)
}

func Deadlocks(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_deadlock",
		Help:        "Number of deadlocks in the last 2m",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			var result float64
			if err := db.QueryRow(`
				SELECT count(*) FROM pg_locks bl
				JOIN pg_stat_activity a
				ON a.pid = bl.pid JOIN pg_locks kl
				ON kl.transactionid = bl.transactionid
				AND kl.pid != bl.pid JOIN pg_stat_activity ka
				ON ka.pid = kl.pid WHERE NOT bl.granted
				AND (ka.query_start >= (now() - interval '2 minutes')
				OR a.query_start >= (now() - interval '2 minutes'))
				AND a.datname = current_database()
			`).Scan(&result); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
			}
			return result / (1024 * 1024)
		},
	)
}
