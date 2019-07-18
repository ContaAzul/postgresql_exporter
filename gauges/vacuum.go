package gauges

import (
	"time"

	"github.com/ContaAzul/postgresql_exporter/postgres"
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

// UnvacuumedTransactions returns the number of unvacuumed transactions
func (g *Gauges) UnvacuumedTransactions() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_unvacuumed_transactions_total",
			Help:        "Number of unvacuumed transactions",
			ConstLabels: g.labels,
		},
		"SELECT age(datfrozenxid) FROM pg_database WHERE datname = current_database()",
	)
}

type tableVacuum struct {
	Name           string  `db:"relname"`
	LastVacuumTime float64 `db:"last_vacuum_time"`
}

// LastTimeVacuumRan returns the last time in seconds at which a table
// was manually vacuumed (not counting VACUUM FULL)
func (g *Gauges) LastTimeVacuumRan() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_last_vacuum_seconds",
			Help:        "Last time in seconds at which a table was manually vacuumed (not counting VACUUM FULL)",
			ConstLabels: g.labels,
		},
		[]string{"table"},
	)

	const lastVacuumQuery = `
		SELECT
		  relname,
		  COALESCE(EXTRACT(EPOCH FROM last_vacuum), 0) as last_vacuum_time
		FROM pg_stat_user_tables
	`

	go func() {
		for {
			var tables []tableVacuum
			if err := g.query(lastVacuumQuery, &tables, emptyParams); err == nil {
				for _, table := range tables {
					gauge.With(prometheus.Labels{
						"table": table.Name,
					}).Set(table.LastVacuumTime)
				}
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}

// LastTimeAutoVacuumRan returns the last time in seconds at which a table
// was vacuumed by the autovacuum daemon
func (g *Gauges) LastTimeAutoVacuumRan() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_last_autovacuum_seconds",
			Help:        "Last time in seconds at which a table was vacuumed by the autovacuum daemon",
			ConstLabels: g.labels,
		},
		[]string{"table"},
	)

	const lastAutoVacuumQuery = `
		SELECT
		  relname,
		  COALESCE(EXTRACT(EPOCH FROM last_autovacuum), 0) as last_vacuum_time
		FROM pg_stat_user_tables
	`

	go func() {
		for {
			var tables []tableVacuum
			if err := g.query(lastAutoVacuumQuery, &tables, emptyParams); err == nil {
				for _, table := range tables {
					gauge.With(prometheus.Labels{
						"table": table.Name,
					}).Set(table.LastVacuumTime)
				}
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}

// VacuumRunningTotal returns the number of backends (including autovacuum worker processes)
// that are currently running a vacuuming (not including VACUUM FULL).
//
// This metric is only supported for PostgreSQL 9.6 or newer versions
func (g *Gauges) VacuumRunningTotal() prometheus.Gauge {
	var gaugeOpts = prometheus.GaugeOpts{
		Name:        "postgresql_vacuum_running_total",
		Help:        "Number of backends (including autovacuum worker processes) currently vacuuming",
		ConstLabels: g.labels,
	}

	const vacuumRunningQuery = `
		SELECT COUNT(*) FROM pg_stat_progress_vacuum WHERE datname = current_database()
	`

	if !postgres.Version(g.version()).IsEqualOrGreaterThan96() {
		log.WithField("db", g.name).
			Warn("postgresql_vacuum_running_total disabled because it's only supported for PostgreSQL 9.6 or newer versions")
		return prometheus.NewGauge(gaugeOpts)
	}
	return g.new(gaugeOpts, vacuumRunningQuery)
}
