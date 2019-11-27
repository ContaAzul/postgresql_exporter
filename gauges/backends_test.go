package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectedBackends(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.ConnectedBackends())
	assert.Len(metrics, 1)
	assertGreaterThan(t, 0, metrics[0])
	assertNoErrs(t, gauges)
}

func TestMaxBackends(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.MaxBackends())
	assert.Len(metrics, 1)
	assertGreaterThan(t, 0, metrics[0])
	assertNoErrs(t, gauges)
}

func TestInstanceConnectedBackends(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.MaxBackends())
	assert.Len(metrics, 1)
	assertGreaterThan(t, 0, metrics[0])
	assertNoErrs(t, gauges)
}

func TestBackendsByState(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.BackendsByState())
	assert.True(len(metrics) > 0)
	for _, m := range metrics {
		assertGreaterThan(t, 0, m)
	}
	assertNoErrs(t, gauges)
}

func TestBackendsByUserAndClientAddress(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.BackendsByUserAndClientAddress())
	assert.True(len(metrics) > 0)
	for _, m := range metrics {
		assertGreaterThan(t, 0, m)
	}
	assertNoErrs(t, gauges)
}

// TODO: somehow set a waiting connections to proper test this
func TestBackendsByWaitEventType(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.BackendsByWaitEventType())
	assert.True(len(metrics) >= 0)
	for _, m := range metrics {
		assertGreaterThan(t, 0, m)
	}
	assertNoErrs(t, gauges)
}
