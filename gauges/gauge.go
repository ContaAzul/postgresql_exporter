package gauges

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

type Gauges struct {
	name     string
	db       *sqlx.DB
	interval time.Duration
	labels   prometheus.Labels
	Errs     prometheus.Gauge
}

func New(name string, db *sql.DB, interval time.Duration) *Gauges {
	var labels = prometheus.Labels{
		"database_name": name,
	}
	return &Gauges{
		name:     name,
		db:       sqlx.NewDb(db, "postgres"),
		interval: interval,
		labels:   labels,
		Errs: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:        "postgresql_query_errors",
				Help:        "queries that failed on the monitoring databases",
				ConstLabels: labels,
			},
		),
	}
}

func paramsFix(params []string) []interface{} {
	iparams := make([]interface{}, len(params))
	for i, v := range params {
		iparams[i] = v
	}
	return iparams
}

func (g *Gauges) new(opts prometheus.GaugeOpts, query string, params ...string) prometheus.Gauge {
	var gauge = prometheus.NewGauge(opts)
	go g.observe(gauge, query, paramsFix(params))
	return gauge
}

func (g *Gauges) fromOnce(gauge prometheus.Gauge, query string, params ...string) {
	go g.observeOnce(gauge, query, paramsFix(params))
}

func (g *Gauges) observeOnce(gauge prometheus.Gauge, query string, params []interface{}) {
	var log = log.WithField("db", g.name)
	log.Debugf("collecting")
	var result []float64
	if err := g.query(query, &result, params); err == nil {
		gauge.Set(result[0])
	}
}

func (g *Gauges) observe(gauge prometheus.Gauge, query string, params []interface{}) {
	for {
		g.observeOnce(gauge, query, params)
		time.Sleep(g.interval)
	}
}

var emptyParams = []interface{}{}

func (g *Gauges) query(query string, result interface{}, params []interface{}) error {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(5*time.Second),
	)
	defer func() {
		<-ctx.Done()
	}()
	var err = g.db.SelectContext(ctx, result, query, params...)
	if err != nil {
		var q = strings.Join(strings.Fields(query), " ")
		g.Errs.Inc()
		log.WithError(err).WithField("query", q).Error("query failed")
	}
	cancel()
	return err
}

var versionRE = regexp.MustCompile(`^PostgreSQL (\d\.\d\.\d).*`)

func (g *Gauges) version() string {
	var version string
	if err := g.db.QueryRow("select version()").Scan(&version); err != nil {
		log.WithError(err).Error("failed to get postgresql version, assuming 9.6.0")
		return "9.6.0"
	}
	return versionRE.FindStringSubmatch(version)[1]
}

func isPG96(version string) bool {
	return strings.HasPrefix(version, "9.6.")
}
