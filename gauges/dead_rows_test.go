package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableDeadRows(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.TableDeadRows())
	assert.True(len(metrics) > 0)
	assertNoErrs(t, gauges)
}

func TestDatabaseDeadRows(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.DatabaseDeadRows())
	assert.True(len(metrics) > 0)
	assertNoErrs(t, gauges)
}
