package gauges

import (
	"database/sql"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

func UnusedIndexes(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_db_unused_indexes",
		Help:        "Dabatase unused indexes count",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			var result float64
			if err := db.QueryRow(`
				SELECT COUNT(*)
				FROM pg_stat_user_indexes ui
				JOIN pg_index i ON ui.indexrelid = i.indexrelid
				WHERE NOT i.indisunique
				AND ui.idx_scan < 100
			`).Scan(&result); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
			}
			return result
		},
	)
}
