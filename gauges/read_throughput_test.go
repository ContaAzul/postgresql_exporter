package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseReadingUsage(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.DatabaseReadingUsage())
	assert.True(len(metrics) > 0)
	assertNoErrs(t, gauges)
}
