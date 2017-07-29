package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplicationStatus(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, ReplicationStatus(db, labels))
	assert.Len(metrics, 1)
	assertGreaterThan(t, -2, metrics[0])
}

func TestReplicationLag(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, ReplicationLag(db, labels))
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}
