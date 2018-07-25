package gauges

import (
	"time"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

type relation struct {
	Name string `db:"relname"`
}

// DeadTuples returns the percentage of dead tuples on the top 20 biggest tables
func (g *Gauges) DeadTuples() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "postgresql_dead_tuples_pct",
		Help:        "percentage of dead tuples on the top 20 biggest tables",
		ConstLabels: g.labels,
	}, []string{"table"})

	if !g.hasExtension("pgstattuple") {
		log.Warn("postgresql_dead_tuples_pct disabled because pgstattuple extension is not installed")
		return gauge
	}
	if !g.hasPermissionToExecutePgStatTuple() {
		log.Warn("postgresql_dead_tuples_pct disabled because user doesn't have permission to use pgstattuple functions")
		return gauge
	}

	const relationsQuery = `
		SELECT relname FROM pg_stat_user_tables ORDER BY n_tup_ins + n_tup_upd desc LIMIT 20
	`

	go func() {
		for {
			var tables []relation
			g.query(relationsQuery, &tables, emptyParams)
			for _, table := range tables {
				var pct []float64
				if err := g.queryWithTimeout(
					"SELECT dead_tuple_percent FROM pgstattuple($1)",
					&pct,
					[]interface{}{table.Name},
					1*time.Minute,
				); err == nil {
					gauge.With(prometheus.Labels{"table": table.Name}).Set(pct[0])
				}
			}
			time.Sleep(12 * time.Hour)
		}
	}()

	return gauge
}

func (g *Gauges) hasPermissionToExecutePgStatTuple() bool {
	if _, err := g.db.Exec("SELECT 1 FROM pgstattuple('pg_class')"); err != nil {
		log.WithError(err).Error("failed to execute pgstattuple function")
		return false
	}
	return true
}
