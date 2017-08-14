package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapBlocksRead(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.HeapBlocksRead())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
	assertNoErrs(t, gauges)
}

func TestHeapBlocksHit(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.HeapBlocksHit())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
	assertNoErrs(t, gauges)
}
