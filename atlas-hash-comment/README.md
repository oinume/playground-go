# What is this?

This is a directory for reproducing the incompatible changes introduced in atlas v0.29.0.

# How to reproduce

In atlas v0.28.1 or previous versions, `atlas migrate apply` succeeds.

```
$ export ATLAS_VERSION=v0.28.1
$ curl -sSf https://atlasgo.sh | sh

$ atlas version

$ atlas migrate apply \
  --dir "file://migrations" \
  --url "postgres://postgres:pass@localhost:5432/atlas-db?search_path=public&sslmode=disable"
```


However, v0.29.0 or later, it fails.

```
$ export ATLAS_VERSION=v0.29.0 
$ curl -sSf https://atlasgo.sh | sh

$ atlas version
atlas version v0.29.0
https://github.com/ariga/atlas/releases/tag/v0.29.0

$ atlas migrate apply \
  --dir "file://migrations" \
  --url "postgres://postgres:pass@localhost:5432/atlas-db?search_path=public&sslmode=disable" 

Migrating to version 20250227054634 (1 migrations in total):

  -- migrating version 20250227054634
    -> # Hash comment
       CREATE TABLE users (
           id text PRIMARY KEY,
           created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
           updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
       );
    pq: syntax error at or near "#"

  -------------------------
  -- 21.258916ms
  -- 1 migration with errors
  -- 1 sql statement with errors
Error: sql/migrate: executing statement "# Hash comment\nCREATE TABLE users (\n    id text PRIMARY KEY,\n    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),\n    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()\n);" from version "20250227054634": pq: syntax error at or near "#"
sql/migrate: write revision: pq: current transaction is aborted, commands ignored until end of transaction block
```