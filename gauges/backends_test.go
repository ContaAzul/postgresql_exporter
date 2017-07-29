package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackends(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, Backends(db, labels))
	assert.Len(metrics, 1)
	assertGreaterThan(t, 0, metrics[0])
}

func TestMaxBackends(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, MaxBackends(db, labels))
	assert.Len(metrics, 1)
	assertGreaterThan(t, 0, metrics[0])
}

func TestWaitingBackends(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, WaitingBackends(db, labels, "9.6.1"))
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}

func TestBackendsStatus(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, BackendsStatus(db, labels)...)
	assert.Len(metrics, 3)
	for _, m := range metrics {
		assertGreaterThan(t, -1, m)
	}
}
