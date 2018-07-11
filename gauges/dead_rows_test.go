package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableDeadRows(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS deadrowstable(id bigint)")
	assert.NoError(err)
	defer close()
	var metrics = evaluate(t, gauges.TableDeadRows())
	assert.True(len(metrics) > 0)
	assertNoErrs(t, gauges)
}

func TestDatabaseDeadRows(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS deadrowstable(id bigint)")
	assert.NoError(err)
	defer close()
	var metrics = evaluate(t, gauges.DatabaseDeadRows())
	assert.True(len(metrics) > 0)
	assertNoErrs(t, gauges)
}
