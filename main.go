package main

import (
	"database/sql"
	"flag"
	"net/http"
	"time"

	"github.com/ContaAzul/postgresql_exporter/config"
	"github.com/ContaAzul/postgresql_exporter/gauges"
	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr       = flag.String("listen-address", ":9111", "The address to listen on for HTTP requests.")
	configFile = flag.String("config", "config.yml", "The path to the config file.")
	interval   = flag.Duration("interval", 20*time.Second, "interval between gathering metrics")
	maxDBConns = flag.Int("max-db-connections", 1, "max connections to open to each database")
	debug      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	log.SetHandler(text.Default)
	flag.Parse()
	var server = &http.Server{
		Addr:         *addr,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}
	if *debug {
		log.SetLevel(log.DebugLevel)
		server.Handler = httplog.New(http.DefaultServeMux)
	}
	http.Handle("/metrics", promhttp.Handler())
	var config = config.Parse(*configFile)
	for _, con := range config.Databases {
		var log = log.WithField("db", con.Name)
		log.Info("started monitoring")
		db, err := sql.Open("postgres", con.URL)
		if err != nil {
			log.WithError(err).Error("failed to open url")
		}
		if err := db.Ping(); err != nil {
			log.WithError(err).Error("failed to ping database")
		}
		db.SetMaxOpenConns(*maxDBConns)
		defer db.Close()

		watch(db, prometheus.DefaultRegisterer, con.Name)
	}
	log.WithField("addr", *addr).Info("started")
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to start up server")
	}
}

func watch(db *sql.DB, reg prometheus.Registerer, name string) {
	var gauges = gauges.New(name, db, *interval)
	reg.MustRegister(gauges.Errs)

	reg.MustRegister(gauges.Backends())
	reg.MustRegister(gauges.BackendsStatus())
	reg.MustRegister(gauges.BufferOversize())
	reg.MustRegister(gauges.BuffersWritten())
	reg.MustRegister(gauges.Deadlocks())
	reg.MustRegister(gauges.DeadTuples())
	reg.MustRegister(gauges.HeapBlocksHit())
	reg.MustRegister(gauges.HeapBlocksRead())
	reg.MustRegister(gauges.IdleSessions())
	reg.MustRegister(gauges.IndexBloat())
	reg.MustRegister(gauges.IndexBlocksHit())
	reg.MustRegister(gauges.IndexBlocksRead())
	reg.MustRegister(gauges.Locks())
	reg.MustRegister(gauges.MaxBackends())
	reg.MustRegister(gauges.ReplicationLag())
	reg.MustRegister(gauges.ReplicationStatus())
	reg.MustRegister(gauges.RequestedCheckpoints())
	reg.MustRegister(gauges.ScheduledCheckpoints())
	reg.MustRegister(gauges.Size())
	reg.MustRegister(gauges.StreamingWALs())
	reg.MustRegister(gauges.TempFiles())
	reg.MustRegister(gauges.TempSize())
	reg.MustRegister(gauges.TransactionsSum())
	reg.MustRegister(gauges.UnusedIndexes())
	reg.MustRegister(gauges.Up())
}
