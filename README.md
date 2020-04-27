# postgresql_exporter

A Prometheus exporter for some postgresql metrics.

## Getting Started

You can add as many database connections as you like to the
`config.yml` file, and run it with:

```console
./postgresql_exporter -config=my/config.yml
```

Then you can add hostname:9111 to the prometheus scrapes config:

```yml
- job_name: 'postgresql'
  static_configs:
    - targets: ['localhost:9111']
```

And voilá, metrics should be there and you should be able to query,
graph and alert on them.

## Setting up a restricted monitoring user

By default some stat views like pg_stat_statements and pg_stat_activity doesn't allow viewing queries run by other users, unless you are a database superuser. Since you probably don't want monitoring to run as a superuser, you can setup a separate monitoring user like this:

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
CREATE EXTENSION IF NOT EXISTS pgstattuple;

CREATE SCHEMA monitoring;

CREATE OR REPLACE FUNCTION monitoring.pgstattuple(IN relname text,
    OUT table_len BIGINT,
    OUT tuple_count BIGINT,
    OUT tuple_len BIGINT,
    OUT tuple_percent FLOAT8,
    OUT dead_tuple_count BIGINT,
    OUT dead_tuple_len BIGINT,
    OUT dead_tuple_percent FLOAT8,
    OUT free_space BIGINT,
    OUT free_percent FLOAT8) AS $$
  SELECT
    table_len,
    tuple_count,
    tuple_len,
    tuple_percent,
    dead_tuple_count,
    dead_tuple_len,
    dead_tuple_percent,
    free_space,
    free_percent
  FROM public.pgstattuple(relname)
$$ LANGUAGE SQL VOLATILE SECURITY DEFINER;

CREATE OR REPLACE FUNCTION monitoring.pgstattuple_approx(IN relname text,
    OUT table_len            BIGINT,
    OUT scanned_percent 	 FLOAT8,
    OUT approx_tuple_count 	 BIGINT,
    OUT approx_tuple_len 	 BIGINT,
    OUT approx_tuple_percent FLOAT8,
    OUT dead_tuple_count 	 BIGINT,
    OUT dead_tuple_len 	     BIGINT,
    OUT dead_tuple_percent 	 FLOAT8,
    OUT approx_free_space 	 BIGINT,
    OUT approx_free_percent  FLOAT8) AS $$
  SELECT
    table_len,
    scanned_percent,
    approx_tuple_count,
    approx_tuple_len,
    approx_tuple_percent,
    dead_tuple_count,
    dead_tuple_len,
    dead_tuple_percent,
    approx_free_space,
    approx_free_percent
  FROM public.pgstattuple_approx(relname)
$$ LANGUAGE SQL VOLATILE SECURITY DEFINER;

CREATE ROLE monitoring WITH LOGIN PASSWORD 'mypassword' 
  CONNECTION LIMIT 5 IN ROLE pg_monitor;
ALTER ROLE monitoring SET search_path = monitoring, pg_catalog, public;

GRANT CONNECT ON DATABASE {{database_name}} TO monitoring;
GRANT USAGE ON SCHEMA monitoring TO monitoring;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA monitoring TO monitoring;
```

Note that these statements must be run as a superuser (to create the SECURITY DEFINER function), but from here onwards you can use the `monitoring` user instead. The exporter will automatically use the helper methods if they exist in the `monitoring` schema, otherwise data will be fetched directly.

The default role `pg_monitor` only has in PostgreSQL 10 or later (See more details [here](https://www.postgresql.org/docs/10/static/default-roles.html)). If you're running Postgres 9.6 or lower you need to create some other helper methods in the `monitoring` schema:

```sql
CREATE OR REPLACE FUNCTION monitoring.pg_stat_activity() RETURNS SETOF pg_stat_activity AS $$
  SELECT * FROM pg_catalog.pg_stat_activity;
$$ LANGUAGE sql VOLATILE SECURITY DEFINER;

CREATE VIEW monitoring.pg_stat_activity AS 
  SELECT * FROM monitoring.pg_stat_activity();

CREATE OR REPLACE FUNCTION monitoring.pg_stat_statements() RETURNS SETOF pg_stat_statements AS $$
  SELECT * FROM public.pg_stat_statements;
$$ LANGUAGE sql VOLATILE SECURITY DEFINER;

CREATE VIEW monitoring.pg_stat_statements AS 
  SELECT * FROM monitoring.pg_stat_statements();

CREATE OR REPLACE FUNCTION monitoring.pg_stat_replication() RETURNS SETOF pg_stat_replication AS $$
  SELECT * FROM pg_catalog.pg_stat_replication;
$$ LANGUAGE sql VOLATILE SECURITY DEFINER;

CREATE VIEW monitoring.pg_stat_replication AS 
  SELECT * FROM monitoring.pg_stat_replication();

CREATE OR REPLACE FUNCTION monitoring.pg_stat_progress_vacuum() RETURNS SETOF pg_stat_progress_vacuum AS $$
  SELECT * FROM pg_catalog.pg_stat_progress_vacuum;
$$ LANGUAGE sql VOLATILE SECURITY DEFINER;

CREATE VIEW monitoring.pg_stat_progress_vacuum AS 
  SELECT * FROM monitoring.pg_stat_progress_vacuum();
```

## Running it within Docker

```console
docker run -p 9111 -v /path/to/my/config.yml:/config.yml caninjas/postgresql_exporter
```