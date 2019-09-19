package gauges

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplicationSlotStatus(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestLogicalReplicationSlot := createTestLogicalReplicationSlot(t, db)
	defer dropTestLogicalReplicationSlot()
	var metrics = evaluate(t, gauges.ReplicationSlotStatus())
	assert.Len(metrics, 1)
	assertEqual(t, 0, metrics[0])
	assertNoErrs(t, gauges)
}

func TestReplicationSlotLagInMegabytes(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestLogicalReplicationSlot := createTestLogicalReplicationSlot(t, db)
	defer dropTestLogicalReplicationSlot()
	var metrics = evaluate(t, gauges.ReplicationSlotLagInMegabytes())
	assert.Len(metrics, 1)
	assertGreaterThan(t, -1, metrics[0])
	assertNoErrs(t, gauges)
}

func createTestLogicalReplicationSlot(t *testing.T, db *sql.DB) func() {
	_, err := db.Exec("SELECT * FROM pg_create_logical_replication_slot('integration_test', 'test_decoding');")
	require.NoError(t, err)
	return func() {
		_, err := db.Exec("SELECT pg_drop_replication_slot('integration_test');")
		assert.New(t).NoError(err)
	}
}
