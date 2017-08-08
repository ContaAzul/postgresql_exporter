package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplicationStatus(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.ReplicationStatus())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -2, metrics[0])
}

func TestReplicationLag(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.ReplicationLag())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}

func TestStreamingWALs(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.StreamingWALs())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}
