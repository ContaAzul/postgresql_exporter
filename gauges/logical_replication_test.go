package gauges

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplicationSlotStatus(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	dropTestLogicalReplicationSlot := createTestLogicalReplicationSlot("test_status", t, db)
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
	dropTestLogicalReplicationSlot := createTestLogicalReplicationSlot("test_lag", t, db)
	defer dropTestLogicalReplicationSlot()
	var metrics = evaluate(t, gauges.ReplicationSlotLagInMegabytes())
	assert.Len(metrics, 0)
	assertNoErrs(t, gauges)
}

func createTestLogicalReplicationSlot(slotName string, t *testing.T, db *sql.DB) func() {
	_, err := db.Exec(fmt.Sprintf("SELECT * FROM pg_create_logical_replication_slot('%s', 'test_decoding');", slotName))
	require.NoError(t, err)
	return func() {
		_, err := db.Exec(fmt.Sprintf("SELECT pg_drop_replication_slot('%s');", slotName))
		assert.New(t).NoError(err)
	}
}
