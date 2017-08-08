package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnusedIndexes(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.UnusedIndexes())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}

func TestIndexBlocksRead(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.IndexBlocksRead())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}

func TestIndexBlocksHit(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.IndexBlocksHit())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}
