package gauges

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

func UnusedIndexes(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	return newGauge(
		db,
		prometheus.GaugeOpts{
			Name:        "postgresql_db_unused_indexes",
			Help:        "Dabatase unused indexes count",
			ConstLabels: labels,
		},
		`
			SELECT COUNT(*)
			FROM pg_stat_user_indexes ui
			JOIN pg_index i ON ui.indexrelid = i.indexrelid
			WHERE NOT i.indisunique
			AND ui.idx_scan < 100
		`,
	)
}
