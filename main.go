package main

import (
	"database/sql"
	"flag"
	"net/http"
	"regexp"
	"time"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/caarlos0/postgresql_exporter/config"
	"github.com/caarlos0/postgresql_exporter/gauges"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr       = flag.String("listen-address", ":9111", "The address to listen on for HTTP requests.")
	configFile = flag.String("config", "config.yml", "The path to the config file.")
)

func main() {
	log.SetHandler(logfmt.Default)
	flag.Parse()
	var config = config.Parse(*configFile)
	for _, con := range config.Databases {
		db, err := sql.Open("postgres", con.URL)
		if err != nil {
			log.WithError(err).Fatal("failed to connect to the database")
		}
		if err := db.Ping(); err != nil {
			log.WithError(err).Fatal("failed to ping the database")
		}
		db.SetMaxOpenConns(1)
		defer db.Close()

		var version = pgVersion(db)
		var labels = prometheus.Labels{
			"database_name": con.Name,
		}

		prometheus.MustRegister(gauges.Up(db, labels, version))
		prometheus.MustRegister(gauges.Size(db, labels))
		prometheus.MustRegister(gauges.IdleSessions(db, labels))
		prometheus.MustRegister(gauges.Backends(db, labels))
		for _, collector := range gauges.BackendsStatus(db, labels) {
			prometheus.MustRegister(collector)
		}
		prometheus.MustRegister(gauges.MaxBackends(db, labels))
		prometheus.MustRegister(gauges.WaitingBackends(db, labels, version))
		prometheus.MustRegister(gauges.UnusedIndexes(db, labels))
		prometheus.MustRegister(gauges.Locks(db, labels))
		prometheus.MustRegister(gauges.ReplicationLag(db, labels))
		prometheus.MustRegister(gauges.Deadlocks(db, labels))
	}

	http.Handle("/metrics", promhttp.Handler())

	var server = &http.Server{
		Handler:      httplog.New(http.DefaultServeMux),
		Addr:         *addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.WithField("addr", *addr).Info("started")
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to start up server")
	}
}

var versionRE = regexp.MustCompile(`^PostgreSQL (\d\.\d\.\d).*`)

func pgVersion(db *sql.DB) string {
	var version string
	if err := db.QueryRow("select version()").Scan(&version); err != nil {
		log.WithError(err).Fatal("failed to get postgresql version")
	}
	return versionRE.FindStringSubmatch(version)[1]
}
