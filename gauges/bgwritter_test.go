package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestedCheckpoints(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.RequestedCheckpoints())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
	assertNoErrs(t, gauges)
}

func TestScheduledCheckpoints(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.ScheduledCheckpoints())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
	assertNoErrs(t, gauges)
}

func TestBufferOversize(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.BufferOversize())
	assert.Len(metrics, 1)
	assert.Equal(float64(0), metrics[0].Value)
	assertNoErrs(t, gauges)
}

func TestBuffersWritten(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.BuffersWritten())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
	assertNoErrs(t, gauges)
}
