package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: somehow create some bloated table to proper test this
func TestTableBloat(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.TableBloat())
	assert.Len(metrics, 0)
	assertNoErrs(t, gauges)
}

func TestTableUsage(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	var metrics = evaluate(t, gauges.TableUsage())
	assert.True(len(metrics) > 0)
	assertNoErrs(t, gauges)
}
