package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var tableDeadRowsQuery = `
	SELECT relname
		 , coalesce(n_dead_tup, 0) as n_dead_tup
	  FROM pg_stat_user_tables
`

type tableDeadRows struct {
	Table      string  `db:"relname"`
	DeadTuples float64 `db:"n_dead_tup"`
}

func (g *Gauges) TableDeadRows() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_table_dead_rows",
			Help:        "Number of dead rows in table",
			ConstLabels: g.labels,
		},
		[]string{"table"},
	)
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

var databaseDeadRowsQuery = `
	SELECT sum(coalesce(n_dead_tup, 0)) as n_dead_tup
	  FROM pg_stat_user_tables
`

func (g *Gauges) DatabaseDeadRows() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_database_dead_rows",
			Help:        "Number of dead rows in database",
			ConstLabels: g.labels,
		}, databaseDeadRowsQuery)
}
