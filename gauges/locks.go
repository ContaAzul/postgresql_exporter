package gauges

import (
	"database/sql"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

func Locks(db *sql.DB, labels prometheus.Labels) prometheus.GaugeFunc {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_db_locks",
		Help:        "Dabatase lock count",
		ConstLabels: labels,
	}
	return prometheus.NewGaugeFunc(
		opts,
		func() float64 {
			var result float64
			if err := db.QueryRow(`
				SELECT count(*)
				FROM pg_locks blocked_locks
				JOIN pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
				JOIN pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
				AND blocking_locks.DATABASE IS NOT DISTINCT
				FROM blocked_locks.DATABASE
				AND blocking_locks.relation IS NOT DISTINCT
				FROM blocked_locks.relation
				AND blocking_locks.page IS NOT DISTINCT
				FROM blocked_locks.page
				AND blocking_locks.tuple IS NOT DISTINCT
				FROM blocked_locks.tuple
				AND blocking_locks.virtualxid IS NOT DISTINCT
				FROM blocked_locks.virtualxid
				AND blocking_locks.transactionid IS NOT DISTINCT
				FROM blocked_locks.transactionid
				AND blocking_locks.classid IS NOT DISTINCT
				FROM blocked_locks.classid
				AND blocking_locks.objid IS NOT DISTINCT
				FROM blocked_locks.objid
				AND blocking_locks.objsubid IS NOT DISTINCT
				FROM blocked_locks.objsubid
				AND blocking_locks.pid != blocked_locks.pid
				JOIN pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
				WHERE NOT blocked_locks.GRANTED
			`).Scan(&result); err != nil {
				log.WithError(err).Warnf("%s: failed to query metric", opts.Name)
			}
			return result
		},
	)
}
