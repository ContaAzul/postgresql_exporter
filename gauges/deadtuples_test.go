package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeadTuples(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS testtable(id bigint)")
	assert.NoError(err)
	var metrics = evaluate(t, gauges.DeadTuples())
	assert.Len(metrics, 1)
	for _, m := range metrics {
		assert.Equal(0.0, m.Value)
	}
}
