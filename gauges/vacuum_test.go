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
