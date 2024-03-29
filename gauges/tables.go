package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const tableBloatQuery = `
WITH constants AS (
	SELECT current_setting('block_size')::numeric AS bs, 23 AS hdr, 8 AS ma
	),
	no_stats AS (
	SELECT table_schema, table_name,
		n_live_tup::numeric as est_rows,
		pg_table_size(relid)::numeric as table_size
	FROM information_schema.columns
		JOIN pg_stat_user_tables as psut
		ON table_schema = psut.schemaname
		AND table_name = psut.relname
		LEFT OUTER JOIN pg_stats
		ON table_schema = pg_stats.schemaname
		AND table_name = pg_stats.tablename
		AND column_name = attname
	WHERE attname IS NULL
	AND table_schema NOT IN ('pg_catalog', 'information_schema')
	GROUP BY table_schema, table_name, relid, n_live_tup
	),
	null_headers AS (
	-- calculate null header sizes
	-- omitting tables which dont have complete stats
	-- and attributes which arent visible
	SELECT
		hdr+1+(sum(case when null_frac <> 0 THEN 1 else 0 END)/8) as nullhdr,
		SUM((1-null_frac)*avg_width) as datawidth,
		MAX(null_frac) as maxfracsum,
		schemaname,
		tablename,
		hdr, ma, bs
	FROM pg_stats CROSS JOIN constants
	LEFT OUTER JOIN no_stats
	ON schemaname = no_stats.table_schema
	AND tablename = no_stats.table_name
	WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
	AND no_stats.table_name IS NULL
	AND EXISTS ( SELECT 1
		FROM information_schema.columns
		WHERE schemaname = columns.table_schema
		AND tablename = columns.table_name )
	   GROUP BY schemaname, tablename, hdr, ma, bs
	),
	data_headers AS (
	-- estimate header and row size
		SELECT
		ma, bs, hdr, schemaname, tablename,
		(datawidth+(hdr+ma-(case when hdr%ma=0 THEN ma ELSE hdr%ma END)))::numeric AS datahdr,
		(maxfracsum*(nullhdr+ma-(case when nullhdr%ma=0 THEN ma ELSE nullhdr%ma END))) AS nullhdr2
		FROM null_headers
	),
	table_estimates AS (
	-- make estimates of how large the table should be
	-- based on row and page size
	SELECT schemaname, tablename, bs,
		reltuples::numeric as est_rows, relpages * bs as table_bytes,
	CEIL((reltuples*
		(datahdr + nullhdr2 + 4 + ma -
		(CASE WHEN datahdr%ma=0 THEN ma ELSE datahdr%ma END)
		)/(bs-20))) * bs AS expected_bytes,
	reltoastrelid
	FROM data_headers
	JOIN pg_class ON tablename = relname
	JOIN pg_namespace ON relnamespace = pg_namespace.oid
	AND schemaname = nspname
	WHERE pg_class.relkind = 'r'
	),
	estimates_with_toast AS (
	SELECT schemaname, tablename,
		TRUE as can_estimate,
		est_rows,
		table_bytes + ( coalesce(toast.relpages, 0) * bs ) as table_bytes,
		expected_bytes + ( ceil( coalesce(toast.reltuples, 0) / 4 ) * bs ) as expected_bytes
	FROM table_estimates LEFT OUTER JOIN pg_class as toast
	ON table_estimates.reltoastrelid = toast.oid
	AND toast.relkind = 't'
	),
	table_estimates_plus AS (
	SELECT current_database() as databasename,
		schemaname, tablename, can_estimate,
		est_rows,
	CASE WHEN table_bytes > 0
		THEN table_bytes::NUMERIC
		ELSE NULL::NUMERIC END
		AS table_bytes,
	CASE WHEN expected_bytes > 0
		THEN expected_bytes::NUMERIC
		ELSE NULL::NUMERIC END
		AS expected_bytes,
	  CASE WHEN expected_bytes > 0 AND table_bytes > 0
		AND expected_bytes <= table_bytes
		THEN (table_bytes - expected_bytes)::NUMERIC
		ELSE 0::NUMERIC END AS bloat_bytes
		FROM estimates_with_toast
		UNION ALL
		SELECT current_database() as databasename,
			table_schema, table_name, FALSE,
			est_rows, table_size,
			NULL::NUMERIC, NULL::NUMERIC
		FROM no_stats
	),
	bloat_data AS (
		select current_database() as databasename,
		schemaname, tablename, can_estimate,
		table_bytes, round(table_bytes/(1024^2)::NUMERIC,3) as table_mb,
		expected_bytes, round(expected_bytes/(1024^2)::NUMERIC,3) as expected_mb,
		round(bloat_bytes*100/table_bytes) as pct_bloat,
		round(bloat_bytes/(1024::NUMERIC^2),2) as mb_bloat,
		table_bytes, expected_bytes, est_rows
	FROM table_estimates_plus
	)
	SELECT tablename,
		COALESCE(pct_bloat,0) as pct_bloat
	FROM bloat_data
--	WHERE ( pct_bloat >= 30 AND mb_bloat >= 10 )
--	OR ( pct_bloat >= 20 AND mb_bloat >= 1000 )
	ORDER BY pct_bloat DESC
`

type tableBloat struct {
	Name string  `db:"tablename"`
	Pct  float64 `db:"pct_bloat"`
}

// TableBloat returns the bloat percentage of a table reporting only for tables
// with bloat percentange greater than 30% or greater than 1000mb
func (g *Gauges) TableBloat() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_table_bloat_pct",
			Help:        "bloat percentage of a table. Reports only for tables with a lot of bloat",
			ConstLabels: g.labels,
		},
		[]string{"table"},
	)
	go func() {
		for {
			gauge.Reset()
			var tables []tableBloat
			if err := g.query(tableBloatQuery, &tables, emptyParams); err == nil {
				for _, table := range tables {
					gauge.With(prometheus.Labels{
						"table": table.Name,
					}).Set(table.Pct)
				}
			}
			time.Sleep(1 * time.Hour)
		}
	}()
	return gauge
}

var tableUsageQuery = `
	SELECT  s.relname,
			coalesce(s.seq_tup_read, 0) as seq_tup_read,
			coalesce(s.idx_tup_fetch, 0) as idx_tup_fetch,
			coalesce(s.n_tup_ins, 0) as n_tup_ins,
			coalesce(s.n_tup_upd, 0) as n_tup_upd,
			coalesce(s.n_tup_del, 0) as n_tup_del
	FROM pg_stat_user_tables s
	ORDER BY 2 desc
`

type tableUsage struct {
	Name      string  `db:"relname"`
	SeqReads  float64 `db:"seq_tup_read"`
	IdxFetchs float64 `db:"idx_tup_fetch"`
	Inserts   float64 `db:"n_tup_ins"`
	Updates   float64 `db:"n_tup_upd"`
	Deletes   float64 `db:"n_tup_del"`
}

func (g *Gauges) TableUsage() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_table_usage",
			Help:        "table usage statistics",
			ConstLabels: g.labels,
		},
		[]string{"table", "stat"},
	)
	go func() {
		for {
			var tables []tableUsage
			if err := g.query(tableUsageQuery, &tables, emptyParams); err == nil {
				for _, table := range tables {
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"stat":  "seq_tup_read",
					}).Set(table.SeqReads)
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"stat":  "idx_tup_fetch",
					}).Set(table.IdxFetchs)
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"stat":  "n_tup_ins",
					}).Set(table.Inserts)
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"stat":  "n_tup_upd",
					}).Set(table.Updates)
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"stat":  "n_tup_del",
					}).Set(table.Deletes)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

var tableSecScansQuery = `
	select relname,
		coalesce(seq_scan, 0) as seq_scan,
		coalesce(idx_scan, 0) as idx_scan
	from pg_stat_user_tables
`

type tableScans struct {
	Name    string  `db:"relname"`
	SecScan float64 `db:"seq_scan"`
	IdxScan float64 `db:"idx_scan"`
}

func (g *Gauges) TableScans() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_table_scans",
			Help:        "table scans statistics",
			ConstLabels: g.labels,
		},
		[]string{"table", "scan"},
	)
	go func() {
		for {
			var tables []tableScans
			if err := g.query(tableSecScansQuery, &tables, emptyParams); err == nil {
				for _, table := range tables {
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"scan":  "seq_scan",
					}).Set(table.SecScan)
					gauge.With(prometheus.Labels{
						"table": table.Name,
						"scan":  "idx_scan",
					}).Set(table.IdxScan)
				}
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}

var hotUpdatesQuery = `
SELECT relname
	 , coalesce(n_tup_hot_upd, 0) as n_tup_hot_upd
  FROM pg_stat_user_tables
`

type tableHotUpdates struct {
	Name       string  `db:"relname"`
	HotUpdates float64 `db:"n_tup_hot_upd"`
}

func (g *Gauges) HOTUpdates() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_hot_updates",
			Help:        "Number of rows Heap-Only tuple updated",
			ConstLabels: g.labels,
		},
		[]string{"table"},
	)
	go func() {
		for {
			var tables []tableHotUpdates
			if err := g.query(hotUpdatesQuery, &tables, emptyParams); err == nil {
				for _, table := range tables {
					gauge.With(prometheus.Labels{
						"table": table.Name,
					}).Set(table.HotUpdates)
				}
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}
