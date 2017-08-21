package gauges

import (
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

const slowQueriesQuery = `
	SELECT total_time, query
	FROM pg_stat_statements
	ORDER BY total_time desc limit 10
`

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
		log.Warn("postgresql_slowest_queries disabled because pg_stat_statements extension is not installed")
		return gauge
	}
	go func() {
		for {
			var queries []slowQuery
			if err := g.query(slowQueriesQuery, &queries, emptyParams); err == nil {
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
