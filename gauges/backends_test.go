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

func TestBackendsStatus(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.BackendsStatus())
	assert.True(len(metrics) > 0)
	for _, m := range metrics {
		assertGreaterThan(t, 0, m)
	}
	assertNoErrs(t, gauges)
}
