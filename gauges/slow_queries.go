package gauges

import (
	"fmt"
	"strings"
	"time"

	"github.com/ContaAzul/postgresql_exporter/postgres"
	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

type slowQuery struct {
	Query string  `db:"query"`
	Time  float64 `db:"total_time"`
}

func (g *Gauges) SlowestQueries() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_slowest_queries",
			Help:        "top 10 slowest queries by accumulated time",
			ConstLabels: g.labels,
		},
		[]string{"query"},
	)
	if !g.hasExtension("pg_stat_statements") {
		log.WithField("db", g.name).
			Warn("postgresql_slowest_queries disabled because pg_stat_statements extension is not installed")
		return gauge
	}
	if !g.hasSharedPreloadLibrary("pg_stat_statements") {
		log.WithField("db", g.name).
			Warn("postgresql_slowest_queries disabled because pg_stat_statements is not on shared_preload_libraries")
		return gauge
	}
	go func() {
		for {
			var queries []slowQuery
			if err := g.query(
				fmt.Sprintf(`
					SELECT %[1]s as total_time, query
					FROM pg_stat_statements
					WHERE dbid = (SELECT datid FROM pg_stat_database WHERE datname = current_database())
					ORDER BY %[1]s desc limit 10`,
					postgres.Version(g.version()).PgStatStatementsTotalTimeColumn(),
				),
				&queries,
				emptyParams,
			); err == nil {
				for _, query := range queries {
					gauge.With(prometheus.Labels{
						"query": strings.Join(strings.Fields(query.Query), " "),
					}).Set(query.Time)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}
