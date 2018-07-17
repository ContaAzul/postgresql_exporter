package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnvacuumedTransactions(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.UnvacuumedTransactions())
	assert.Len(metrics, 1)
	assertGreaterThan(t, 0, metrics[0])
	assertNoErrs(t, gauges)
}

func TestLastTimeVacuumRan(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestTable := createTestTable(t, db)
	defer dropTestTable()

	var metrics = evaluate(t, gauges.LastTimeVacuumRan())
	assert.Len(metrics, 1)
	assertEqual(t, 0, metrics[0])
	assertNoErrs(t, gauges)
}

func TestLastTimeAutoVacuumRan(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestTable := createTestTable(t, db)
	defer dropTestTable()

	var metrics = evaluate(t, gauges.LastTimeAutoVacuumRan())
	assert.Len(metrics, 1)
	assertEqual(t, 0, metrics[0])
	assertNoErrs(t, gauges)
}

func TestVacuumRunningTotal(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()

	var metrics = evaluate(t, gauges.VacuumRunningTotal())
	assert.Len(metrics, 1)
	assertEqual(t, 0, metrics[0])
	assertNoErrs(t, gauges)
}
