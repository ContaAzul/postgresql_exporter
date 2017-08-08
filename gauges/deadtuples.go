package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Relation struct {
	Name string `db:"relname"`
}

func (g *Gauges) DeadTuples() *prometheus.GaugeVec {
	var opts = prometheus.GaugeOpts{
		Name:        "postgresql_dead_tuples_pct",
		Help:        "dead tuples percentage on the top 20 biggest tables",
		ConstLabels: g.labels,
	}
	var gauge = prometheus.NewGaugeVec(opts, []string{"table"})

	go func() {
		for {
			var tables []Relation
			g.query(
				`
					SELECT relname
					FROM pg_stat_user_tables
					ORDER BY n_tup_ins + n_tup_upd desc limit 20
				`,
				&tables,
				emptyParams,
			)
			for _, table := range tables {
				g.fromOnce(
					gauge.With(prometheus.Labels{"table": table.Name}),
					"SELECT dead_tuple_percent FROM pgstattuple($1)",
					table.Name,
				)
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}
