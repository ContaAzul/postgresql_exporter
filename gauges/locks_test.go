package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocks(t *testing.T) {
	var assert = assert.New(t)
	var db = connect(t)
	defer db.Close()
	var metrics = evaluate(t, Locks(db, labels))
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
}
