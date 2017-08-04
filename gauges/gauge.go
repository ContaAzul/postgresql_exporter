package gauges

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Gauges struct {
	name     string
	db       *sql.DB
	interval time.Duration
	labels   prometheus.Labels
}

func New(name string, db *sql.DB, interval time.Duration) *Gauges {
	return &Gauges{
		name:     name,
		db:       db,
		interval: interval,
		labels: prometheus.Labels{
			"database_name": name,
		},
	}
}

func (g *Gauges) new(
	opts prometheus.GaugeOpts,
	query string,
	params ...string,
) prometheus.Gauge {
	iparams := make([]interface{}, len(params))
	for i, v := range params {
		iparams[i] = v
	}
	var gauge = prometheus.NewGauge(opts)
	go g.observe(gauge, opts.Name, query, iparams)
	return gauge
}

func (g *Gauges) observe(gauge prometheus.Gauge, metric, query string, params []interface{}) {
	for {
		var result float64
		var log = log.WithField("db", g.name).WithField("metric", metric)
		log.Debugf("collecting")
		ctx, cancel := context.WithDeadline(
			context.Background(),
			time.Now().Add(1*time.Second),
		)
		defer func() {
			<-ctx.Done()
		}()
		var err = g.db.QueryRowContext(ctx, query, params...).Scan(&result)
		if err != nil {
			log.WithError(err).Warn("failed to query metric")
		}
		cancel()
		if err == nil {
			gauge.Set(result)
		}
		time.Sleep(g.interval)
	}
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
