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

func newGauge(
	db *sql.DB,
	opts prometheus.GaugeOpts,
	query string,
	params ...string,
) prometheus.Gauge {
	iparams := make([]interface{}, len(params))
	for i, v := range params {
		iparams[i] = v
	}
	var gauge = prometheus.NewGauge(opts)
	go func() {
		for {
			var log = log.WithField("metric", opts.Name)
			log.Debugf("collecting")
			var result float64
			ctx, cancel := context.WithDeadline(
				context.Background(),
				time.Now().Add(1*time.Second),
			)
			defer func() {
				<-ctx.Done()
			}()
			if err := db.QueryRowContext(ctx, query, iparams...).Scan(&result); err != nil {
				log.WithError(err).Warn("failed to query metric")
			}
			cancel()
			gauge.Set(result)
			time.Sleep(20 * time.Second)
		}
	}()
	return gauge
}

func isPG96(version string) bool {
	return strings.HasPrefix(version, "9.6.")
}

var versionRE = regexp.MustCompile(`^PostgreSQL (\d\.\d\.\d).*`)

func pgVersion(db *sql.DB) string {
	var version string
	if err := db.QueryRow("select version()").Scan(&version); err != nil {
		log.WithError(err).Fatal("failed to get postgresql version")
	}
	return versionRE.FindStringSubmatch(version)[1]
}
