package gauges

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/apex/log"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

type Metric struct {
	Value float64
	Name  string
}

var labels = prometheus.Labels{
	"testing": "true",
}

func evaluate(t *testing.T, gauges ...prometheus.Gauge) (result []Metric) {
	var assert = assert.New(t)
	var reg = prometheus.NewRegistry()
	for _, gauge := range gauges {
		assert.NoError(reg.Register(gauge))
	}
	time.Sleep(100 * time.Millisecond)
	metrics, err := reg.Gather()
	assert.NoError(err)
	for _, metric := range metrics {
		for _, m := range metric.GetMetric() {
			result = append(
				result,
				Metric{
					Value: m.GetGauge().GetValue(),
					Name:  metric.GetName(),
				},
			)
		}
	}
	return
}

func assertGreaterThan(t *testing.T, expected float64, m Metric) {
	var assert = assert.New(t)
	assert.True(
		m.Value > expected,
		"%s should be > %v: %v", m.Name, expected, m.Value,
	)
}

func prepare(t *testing.T) (*Gauges, func()) {
	var db = connect(t)
	var gauges = New("test", db, 1*time.Minute)
	return gauges, func() {
		log.Info("CLOSING")
		assert.NoError(t, db.Close())
	}
}

func connect(t *testing.T) *sql.DB {
	var assert = assert.New(t)
	var url = os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = "postgres://localhost:5432/postgres?sslmode=disable"
	}
	db, err := sql.Open("postgres", url)
	assert.NoError(err, "failed to open connection to the database")
	assert.NoError(db.Ping(), "failed to ping database")
	db.SetMaxOpenConns(1)
	return db
}
