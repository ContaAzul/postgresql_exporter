package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: somehow create some bloated table to proper test this
func TestTableBloat(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestTable := createTestTable(t, db)
	defer dropTestTable()

	var metrics = evaluate(t, gauges.TableBloat())
	assert.Len(metrics, 0)
	assertNoErrs(t, gauges)
}

func TestTableUsage(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestTable := createTestTable(t, db)
	defer dropTestTable()

	var metrics = evaluate(t, gauges.TableUsage())
	assert.True(len(metrics) > 0)
	assertNoErrs(t, gauges)
}

func TestTableScans(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestTable := createTestTable(t, db)
	defer dropTestTable()

	var metrics = evaluate(t, gauges.TableScans())
	assert.True(len(metrics) > 0)
	assertNoErrs(t, gauges)
}

func TestHOTUpdates(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestTable := createTestTable(t, db)
	defer dropTestTable()

	var metrics = evaluate(t, gauges.HOTUpdates())
	assert.True(len(metrics) > 0)
	assertNoErrs(t, gauges)
}
