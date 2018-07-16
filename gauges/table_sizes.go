package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const tableSizesQuery = `
	SELECT
		ut.relname AS table_name,
		pg_indexes_size(c.oid) AS index_size,
		COALESCE(pg_total_relation_size(c.reltoastrelid), 0) AS toast_size,
		pg_total_relation_size(c.oid)
			- pg_indexes_size(c.oid)
			- COALESCE(pg_total_relation_size(c.reltoastrelid), 0) AS table_size
	FROM
		pg_stat_user_tables ut
		JOIN pg_class c on ut.relname = c.relname
`

type tableSizes struct {
	Name      string  `db:"table_name"`
	IndexSize float64 `db:"index_size"`
	ToastSize float64 `db:"toast_size"`
	TableSize float64 `db:"table_size"`
}

// TableSizes returns the total disk space in bytes used by the a table,
// including all indexes and TOAST data
func (g *Gauges) TableSizes() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_table_size_bytes",
			Help:        "Total disk space in bytes used by the a table, including all indexes and TOAST data",
			ConstLabels: g.labels,
		},
		[]string{"table", "type"},
	)

	go func() {
		for {
			var tables []tableSizes
			if err := g.query(tableSizesQuery, &tables, emptyParams); err == nil {
				for _, table := range tables {
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"type":  "index",
					}).Set(table.IndexSize)
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"type":  "toast",
					}).Set(table.ToastSize)
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"type":  "table",
					}).Set(table.TableSize)
				}
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}
