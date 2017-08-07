package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUp(t *testing.T) {
	var assert = assert.New(t)
	gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.Up())
	assert.Len(metrics, 1)
	assert.Equal(1.0, metrics[0].Value, "%s should be 1 ", metrics[0].Name)
}

func TestSize(t *testing.T) {
	var assert = assert.New(t)
	gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.Size())
	assert.Len(metrics, 1)
	assertGreaterThan(t, 1000, metrics[0])
}

func TestDeadlocks(t *testing.T) {
	var assert = assert.New(t)
	gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.Deadlocks())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}

func TestTempSize(t *testing.T) {
	var assert = assert.New(t)
	gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.TempSize())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}
