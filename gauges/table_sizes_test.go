package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableSizes(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestTable := createTestTable(t, db)
	defer dropTestTable()

	var metrics = evaluate(t, gauges.TableSizes())
	assert.Len(metrics, 3)
	for _, metric := range metrics {
		assertGreaterThan(t, -1, metric)
	}
	assertNoErrs(t, gauges)
}
