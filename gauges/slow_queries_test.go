package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlowestQueries(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	if !gauges.hasSharedPreloadLibrary("pg_stat_statements") {
		t.Skip("pg_stat_statements not in shared_preload_libraries")
		return
	}
	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_stat_statements")
	require.NoError(t, err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS slowest_queries AS SELECT generate_series(1, 200) AS id, md5(random()::text) AS desc")
	require.NoError(t, err)
	defer func() {
		_, err := db.Exec("DROP TABLE IF EXISTS slowest_queries")
		assert.NoError(err)
	}()

	for i := 0; i < 20; i++ {
		var id int64
		db.QueryRow("select id from slowest_queries where id = $1", i).Scan(&id)
	}

	var metrics = evaluate(t, gauges.SlowestQueries())
	assert.Len(metrics, 10)
	assertNoErrs(t, gauges)
}

func TestSlowestQueriesExtensionNotInstalled(t *testing.T) {
	var assert = assert.New(t)
	db, gauges, close := prepare(t)
	defer close()
	if !gauges.hasSharedPreloadLibrary("pg_stat_statements") {
		t.Skip("pg_stat_statements not in shared_preload_libraries")
		return
	}
	_, err := db.Exec("DROP EXTENSION IF EXISTS pg_stat_statements")
	require.NoError(t, err)

	var metrics = evaluate(t, gauges.SlowestQueries())
	assert.Len(metrics, 0)
	assertNoErrs(t, gauges)
}
