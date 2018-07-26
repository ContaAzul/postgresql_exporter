package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type tableDeadRows struct {
	Table      string  `db:"relname"`
	DeadTuples float64 `db:"n_dead_tup"`
}

// TableDeadRows returns the estimated number of dead rows of a given table
func (g *Gauges) TableDeadRows() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_table_dead_rows",
			Help:        "Estimated number of dead rows in a table",
			ConstLabels: g.labels,
		},
		[]string{"table"},
	)

	const tableDeadRowsQuery = `
		SELECT relname, coalesce(n_dead_tup, 0) as n_dead_tup FROM pg_stat_user_tables
	`

	go func() {
		for {
			var tableDeadRows []tableDeadRows
			if err := g.query(tableDeadRowsQuery, &tableDeadRows, emptyParams); err == nil {
				for _, table := range tableDeadRows {
					gauge.With(prometheus.Labels{
						"table": table.Table,
					}).Set(table.DeadTuples)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

// DatabaseDeadRows returns the sum of estimated number of dead rows of all tables in a database
func (g *Gauges) DatabaseDeadRows() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_database_dead_rows",
			Help:        "Estimated number of dead rows in a database",
			ConstLabels: g.labels,
		},
		"SELECT coalesce(sum(n_dead_tup), 0) as n_dead_tup FROM pg_stat_user_tables",
	)
}

type tableLiveRows struct {
	Table      string  `db:"relname"`
	LiveTuples float64 `db:"n_live_tup"`
}

// TableLiveRows returns the estimated number of live rows of a given table
func (g *Gauges) TableLiveRows() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_table_live_rows",
			Help:        "Estimated number of live rows in a table",
			ConstLabels: g.labels,
		},
		[]string{"table"},
	)

	const tableLiveRowsQuery = `
		SELECT relname, coalesce(n_live_tup, 0) as n_live_tup FROM pg_stat_user_tables
	`

	go func() {
		for {
			var tableLiveRows []tableLiveRows
			if err := g.query(tableLiveRowsQuery, &tableLiveRows, emptyParams); err == nil {
				for _, table := range tableLiveRows {
					gauge.With(prometheus.Labels{
						"table": table.Table,
					}).Set(table.LiveTuples)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

// DatabaseLiveRows returns the sum of estimated number of live rows of all tables in a database
func (g *Gauges) DatabaseLiveRows() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_database_live_rows",
			Help:        "Estimated number of live rows in a database",
			ConstLabels: g.labels,
		},
		"SELECT coalesce(sum(n_live_tup), 0) as n_live_tup FROM pg_stat_user_tables",
	)
}
