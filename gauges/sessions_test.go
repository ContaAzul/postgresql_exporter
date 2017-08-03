package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdleSessions(t *testing.T) {
	var assert = assert.New(t)
	gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.IdleSessions())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}
