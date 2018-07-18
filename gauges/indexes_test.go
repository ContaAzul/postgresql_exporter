package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexScans(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestTable := createTestTable(t, db)
	defer dropTestTable()

	var metrics = evaluate(t, gauges.IndexScans())
	assert.Len(metrics, 1)
	assertEqual(t, 0, metrics[0])
	assertNoErrs(t, gauges)
}

func TestUnusedIndexes(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.UnusedIndexes())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
	assertNoErrs(t, gauges)
}

func TestIndexBlocksRead(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.IndexBlocksRead())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
	assertNoErrs(t, gauges)
}

func TestIndexBlocksHit(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.IndexBlocksHit())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
	assertNoErrs(t, gauges)
}

// TODO: somehow create some bloated index to proper test this
func TestIndexBloat(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.IndexBloat())
	assert.Len(metrics, 0)
	assertNoErrs(t, gauges)
}
