package main

import (
	"database/sql"
	"flag"
	"net/http"
	"time"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/caarlos0/postgresql_exporter/config"
	"github.com/caarlos0/postgresql_exporter/gauges"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr       = flag.String("listen-address", ":9111", "The address to listen on for HTTP requests.")
	configFile = flag.String("config", "config.yml", "The path to the config file.")
	debug      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	log.SetHandler(text.Default)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	flag.Parse()
	var config = config.Parse(*configFile)
	for _, con := range config.Databases {
		db, err := sql.Open("postgres", con.URL)
		var log = log.WithField("db", con.Name)
		if err != nil {
			log.WithError(err).Fatal("failed to connect to the database")
		}
		if err := db.Ping(); err != nil {
			log.WithError(err).Fatal("failed to ping the database")
		}
		db.SetMaxOpenConns(1)
		defer db.Close()

		watch(db, prometheus.DefaultRegisterer, con.Name)
	}

	http.Handle("/metrics", promhttp.Handler())

	var server = &http.Server{
		Addr:         *addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if *debug {
		server.Handler = httplog.New(http.DefaultServeMux)
	}

	log.WithField("addr", *addr).Info("started")
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to start up server")
	}
}

func watch(db *sql.DB, reg prometheus.Registerer, name string) {
	var labels = prometheus.Labels{
		"database_name": name,
	}
	reg.MustRegister(gauges.Up(db, labels))
	reg.MustRegister(gauges.Size(db, labels))
	reg.MustRegister(gauges.IdleSessions(db, labels))
	reg.MustRegister(gauges.Backends(db, labels))
	for _, collector := range gauges.BackendsStatus(db, labels) {
		reg.MustRegister(collector)
	}
	reg.MustRegister(gauges.MaxBackends(db, labels))
	reg.MustRegister(gauges.WaitingBackends(db, labels))
	reg.MustRegister(gauges.UnusedIndexes(db, labels))
	reg.MustRegister(gauges.Locks(db, labels))
	reg.MustRegister(gauges.ReplicationStatus(db, labels))
	reg.MustRegister(gauges.ReplicationLag(db, labels))
	reg.MustRegister(gauges.Deadlocks(db, labels))
}
