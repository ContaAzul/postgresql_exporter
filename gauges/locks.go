package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type lockCountWithMode struct {
	Mode  string  `db:"mode"`
	Count float64 `db:"count"`
}

// Locks returns the number of locks by mode
func (g *Gauges) Locks() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_lock_count",
			Help:        "Number of locks by mode",
			ConstLabels: g.labels,
		},
		[]string{"mode"},
	)
	go func() {
		for {
			var locks []lockCountWithMode
			if err := g.query(
				`
					SELECT mode, count(*) as count
					FROM pg_locks
					WHERE database = (
						SELECT datid
						FROM pg_stat_database
						WHERE datname = current_database()
					) GROUP BY mode;
				`,
				&locks,
				emptyParams,
			); err == nil {
				for _, lock := range locks {
					gauge.With(prometheus.Labels{
						"mode": lock.Mode,
					}).Set(lock.Count)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

// NotGrantedLocks returns the number of not granted locks
func (g *Gauges) NotGrantedLocks() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_not_granted_locks",
			Help:        "Number of not granted locks",
			ConstLabels: g.labels,
		},
		`
			SELECT count(*) as count
			FROM pg_locks
			WHERE NOT granted
			AND database = (
				SELECT datid
				FROM pg_stat_database
				WHERE datname = current_database()
			);
		`,
	)
}
