package gauges

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

func Backends(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	return newGauge(
		db,
		prometheus.GaugeOpts{
			Name:        "postgresql_totalbackends",
			Help:        "Total database backends",
			ConstLabels: labels,
		},
		`
			SELECT numbackends
			FROM pg_stat_database
			WHERE datname = current_database()
		`,
	)
}

func MaxBackends(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	return newGauge(
		db,
		prometheus.GaugeOpts{
			Name:        "postgresql_max_backends",
			Help:        "Maximum database backends (per postmaster)",
			ConstLabels: labels,
		},
		`
			SELECT setting::numeric
			FROM pg_settings
			WHERE name = 'max_connections'
		`,
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
		result = append(result, newGauge(
			db,
			prometheus.GaugeOpts{
				Name:        "postgresql_backends",
				Help:        "Active database connections",
				ConstLabels: lbl,
			},
			`
				SELECT COUNT(*)
				FROM pg_stat_activity
				WHERE datname = current_database()
				AND state = $1
			`,
			status,
		))
	}
	return result
}

func WaitingBackends(db *sql.DB, labels prometheus.Labels, version string) prometheus.GaugeFunc {
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
	return newGauge(
		db,
		prometheus.GaugeOpts{
			Name:        "postgresql_waiting_backends",
			Help:        "Database connections waiting on a Lock",
			ConstLabels: labels,
		},
		query,
	)
}
