package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type indexScans struct {
	Table   string  `db:"relname"`
	Index   string  `db:"indexrelname"`
	IdxScan float64 `db:"idx_scan"`
}

// IndexScans returns the number of index scans initiated on a index
func (g *Gauges) IndexScans() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_index_scans_total",
			Help:        "Number of index scans initiated on a index",
			ConstLabels: g.labels,
		},
		[]string{"table", "index"},
	)

	const indexScansQuery = "SELECT relname, indexrelname, idx_scan FROM pg_stat_user_indexes"

	go func() {
		for {
			var indexes []indexScans
			if err := g.query(indexScansQuery, &indexes, emptyParams); err == nil {
				for _, index := range indexes {
					gauge.With(prometheus.Labels{
						"table": index.Table,
						"index": index.Index,
					}).Set(index.IdxScan)
				}
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}

// UnusedIndexes returns the count of unused indexes in the database
func (g *Gauges) UnusedIndexes() prometheus.Gauge {
	return g.new(
		prometheus.GaugeOpts{
			Name:        "postgresql_unused_indexes",
			Help:        "Dabatase unused indexes count",
			ConstLabels: g.labels,
		},
		`
			SELECT COUNT(*)
			FROM pg_stat_user_indexes ui
			JOIN pg_index i ON ui.indexrelid = i.indexrelid
			WHERE NOT i.indisunique
			AND ui.idx_scan < 100
		`,
	)
}

type schemaIndexBlocksRead struct {
	Name     		string  `db:"schemaname"`
	IndexBlocksRead	float64 `db:"idx_blks_read"`
}

// IndexBlocksRead returns the sum of the number of disk blocks read from all public indexes
func (g *Gauges) IndexBlocksReadBySchema() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_index_blocks_read_sum",
			Help:        "Sum of the number of disk blocks read from all user indexes",
			ConstLabels: g.labels,
		},
		[]string{"schema"},
	)

	const schemaIndexBlocksReadQuery = `
		SELECT
			schemaname,
			coalesce(sum(idx_blks_read), 0) AS idx_blks_read
		FROM pg_statio_user_indexes
		WHERE schemaname NOT IN ('pg_catalog','information_schema','monitoring')
		GROUP BY schemaname;
	`

	go func() {
		for {
			var schemas []schemaIndexBlocksRead
			if err := g.query(schemaIndexBlocksReadQuery, &schemas, emptyParams); err == nil {
				for _, schema := range schemas {
					gauge.With(prometheus.Labels{
						"schema": schema.Name,
					}).Set(schema.IndexBlocksRead)
				}
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}


type schemaIndexBlocksHit struct {
	Name     		string  `db:"schemaname"`
	IndexBlocksHit	float64 `db:"idx_blks_hit"`
}

// IndexBlocksHit returns the sum of the number of buffer hits on all user indexes
func (g *Gauges) IndexBlocksHitBySchema() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_index_blocks_hit_sum",
			Help:        "Sum of the number of buffer hits on all user indexes",
			ConstLabels: g.labels,
		},
		[]string{"schema"},
	)

	const schemaIndexBlocksHitQuery = `
		SELECT
			schemaname,
			coalesce(sum(idx_blks_hit), 0) AS idx_blks_hit
		FROM pg_statio_user_indexes
		WHERE schemaname NOT IN ('pg_catalog','information_schema','monitoring')
		GROUP BY schemaname;
	`

	go func() {
		for {
			var schemas []schemaIndexBlocksHit
			if err := g.query(schemaIndexBlocksHitQuery, &schemas, emptyParams); err == nil {
				for _, schema := range schemas {
					gauge.With(prometheus.Labels{
						"schema": schema.Name,
					}).Set(schema.IndexBlocksHit)
				}
			}
			time.Sleep(g.interval)
		}
	}()

	return gauge
}

const indexBloatQuery = `
WITH btree_index_atts AS (
	SELECT nspname,
	indexclass.relname as index_name,
	indexclass.reltuples,
	indexclass.relpages,
	indrelid, indexrelid,
	indexclass.relam,
	tableclass.relname as tablename,
	regexp_split_to_table(indkey::text, ' ')::smallint AS attnum,
	indexrelid as index_oid
	FROM pg_index
	JOIN pg_class AS indexclass ON pg_index.indexrelid = indexclass.oid
	JOIN pg_class AS tableclass ON pg_index.indrelid = tableclass.oid
	JOIN pg_namespace ON pg_namespace.oid = indexclass.relnamespace
	JOIN pg_am ON indexclass.relam = pg_am.oid
	WHERE pg_am.amname = 'btree' and indexclass.relpages > 0
	AND nspname NOT IN ('pg_catalog','information_schema')
	),
index_item_sizes AS (
	SELECT
	ind_atts.nspname, ind_atts.index_name,
	ind_atts.reltuples, ind_atts.relpages, ind_atts.relam,
	indrelid AS table_oid, index_oid,
	current_setting('block_size')::numeric AS bs,
	8 AS maxalign,
	24 AS pagehdr,
	CASE WHEN max(coalesce(pg_stats.null_frac,0)) = 0
	THEN 2
	ELSE 6
	END AS index_tuple_hdr,
	sum( (1-coalesce(pg_stats.null_frac, 0)) * coalesce(pg_stats.avg_width, 1024) ) AS nulldatawidth
	FROM pg_attribute
	JOIN btree_index_atts AS ind_atts ON pg_attribute.attrelid = ind_atts.indexrelid AND pg_attribute.attnum = ind_atts.attnum
	JOIN pg_stats ON pg_stats.schemaname = ind_atts.nspname
	-- stats for regular index columns
	AND ( (pg_stats.tablename = ind_atts.tablename AND pg_stats.attname = pg_catalog.pg_get_indexdef(pg_attribute.attrelid, pg_attribute.attnum, TRUE))
	-- stats for functional indexes
	OR   (pg_stats.tablename = ind_atts.index_name AND pg_stats.attname = pg_attribute.attname))
	WHERE pg_attribute.attnum > 0
	GROUP BY 1, 2, 3, 4, 5, 6, 7, 8, 9
),
index_aligned_est AS (
	SELECT maxalign, bs, nspname, index_name, reltuples,
	relpages, relam, table_oid, index_oid,
		coalesce (
		ceil (
			reltuples * ( 6
			+ maxalign
				- CASE
				WHEN index_tuple_hdr%maxalign = 0 THEN maxalign
					ELSE index_tuple_hdr%maxalign
				END
					+ nulldatawidth
					+ maxalign
					- CASE /* Add padding to the data to align on MAXALIGN */
					WHEN nulldatawidth::integer%maxalign = 0 THEN maxalign
					ELSE nulldatawidth::integer%maxalign
					END
			)::numeric
			/ ( bs - pagehdr::NUMERIC )
		+1 )
		, 0 )
	as expected
	FROM index_item_sizes
),
raw_bloat AS (
	SELECT current_database() as dbname, nspname, pg_class.relname AS table_name, index_name,
	bs*(index_aligned_est.relpages)::bigint AS totalbytes, expected,
		CASE
		WHEN index_aligned_est.relpages <= expected
			THEN 0
			ELSE bs*(index_aligned_est.relpages-expected)::bigint
			END AS wastedbytes,
			CASE
			WHEN index_aligned_est.relpages <= expected
				THEN 0
				ELSE bs*(index_aligned_est.relpages-expected)::bigint * 100 / (bs*(index_aligned_est.relpages)::bigint)
			END AS realbloat,
		pg_relation_size(index_aligned_est.table_oid) as table_bytes,
		stat.idx_scan as index_scans
	FROM index_aligned_est
	JOIN pg_class ON pg_class.oid=index_aligned_est.table_oid
	JOIN pg_stat_user_indexes AS stat ON index_aligned_est.index_oid = stat.indexrelid
),
format_bloat AS (
SELECT dbname as database_name, nspname as schema_name, table_name, index_name,
	round(realbloat) as bloat_pct, round(wastedbytes/(1024^2)::NUMERIC) as bloat_mb,
	round(totalbytes/(1024^2)::NUMERIC,3) as index_mb,
	round(table_bytes/(1024^2)::NUMERIC,3) as table_mb,
	index_scans
FROM raw_bloat
)
SELECT table_name, index_name, bloat_pct
FROM format_bloat
WHERE database_name = current_database()
AND bloat_pct > 50
AND bloat_mb > 10
ORDER BY bloat_mb DESC
`

type indexBloat struct {
	Table string  `db:"table_name"`
	Name  string  `db:"index_name"`
	Pct   float64 `db:"bloat_pct"`
}

// IndexBloat returns bloat percentage of an index reporting only for indexes
// with size greater than 10mb and bloat lower than 50%
func (g *Gauges) IndexBloat() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_index_bloat_pct",
			Help:        "Bloat percentage of an index. This metric reports only indexes > 10mb and > 50% bloat",
			ConstLabels: g.labels,
		},
		[]string{"index", "table"},
	)
	go func() {
		for {
			var indexes []indexBloat
			if err := g.query(indexBloatQuery, &indexes, emptyParams); err == nil {
				for _, idx := range indexes {
					gauge.With(prometheus.Labels{
						"table": idx.Table,
						"index": idx.Name,
					}).Set(idx.Pct)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}
