package gauges

import (
	"database/sql"
	"strings"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

func Backends(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_db_numbackends",
		Help:        "Total database backends",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			var result float64
			if err := db.QueryRow(`
				SELECT numbackends
				FROM pg_stat_database
				WHERE datname = current_database()
			`).Scan(&result); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
			}
			return result
		},
	)
}

func MaxBackends(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_db_max_backends",
		Help:        "Maximum database backends (per postmaster)",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			var result float64
			if err := db.QueryRow(`
				SELECT setting::numeric
				FROM pg_settings
				WHERE name = 'max_connections'
			`).Scan(&result); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
			}
			return result
		},
	)
}

func BackendsStatus(db *sql.DB, labels prometheus.Labels) []prometheus.GaugeFunc {
	var result = []prometheus.GaugeFunc{}
	for _, status := range []string{"active", "idle", "idle in transaction"} {
		lbl := prometheus.Labels{}
		for k, v := range labels {
			lbl[k] = v
		}
		lbl["status"] = status
		var opts = prometheus.GaugeOpts{
			Name:        "postgresql_db_backends",
			Help:        "Active database connections",
			ConstLabels: lbl,
		}
		result = append(result, prometheus.NewGaugeFunc(
			opts,
			countBackendsByState(db, opts, status),
		))
	}
	return result
}

func countBackendsByState(db *sql.DB, opts prometheus.GaugeOpts, state string) func() float64 {
	return func() float64 {
		var result float64
		if err := db.QueryRow(`
			SELECT COUNT(*)
			FROM pg_stat_activity
			WHERE datname = current_database()
			AND state = $1
		`, state).Scan(&result); err != nil {
			log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
		}
		return result
	}
}

func WaitingBackends(db *sql.DB, labels prometheus.Labels, version string) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_db_waiting_backends",
		Help:        "Database connections waiting on a Lock",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			var result float64
			var query = `
				SELECT COUNT(*)
				FROM pg_stat_activity
				WHERE datname = current_database()
				AND waiting is true
			`
			if isPG96(version) {
				query = `
					SELECT COUNT(*)
					FROM pg_stat_activity
					WHERE datname = current_database()
					AND wait_event is not null
				`
			}
			if err := db.QueryRow(query).Scan(&result); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
			}
			return result
		},
	)
}

func isPG96(version string) bool {
	return strings.HasPrefix(version, "9.6.")
}
