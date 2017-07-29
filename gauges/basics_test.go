package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUp(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, Up(db, labels))
	assert.Len(metrics, 1)
	assert.Equal(1.0, metrics[0].Value, "%s should be 1 ", metrics[0].Name)
}

func TestSize(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, Size(db, labels))
	assert.Len(metrics, 1)
	assertGreaterThan(t, 1000, metrics[0])
}

func TestDeadlocks(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, Deadlocks(db, labels))
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}
