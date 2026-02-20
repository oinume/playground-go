# PostgreSQL serial vs UUID パフォーマンス検証

PostgreSQL における PK の型として `BIGSERIAL` と `UUID v7` でどの程度パフォーマンスに差があるかを検証する。

## 検証ポイント

1. **INSERT 時のパフォーマンス** — DB 自動採番 vs Go 側での UUID v7 生成
2. **参照の局所性（キャッシュヒット率）** — `shared_buffers=8MB` の制約下で、最近データへのアクセスとランダムアクセスの差

## テーブル定義

| テーブル | PK 型 | 備考 |
|---|---|---|
| `users_serial` | `BIGSERIAL` | DB 自動採番（連番） |
| `users_uuid` | `UUID` | Go 側で `uuid.NewV7()` 生成（時系列ソート可能） |

## セットアップ

### 1. PostgreSQL 起動

```bash
cd postgres_serial_vs_uuid
docker compose up -d
```

- PostgreSQL 18 / ポート 5433
- `shared_buffers=8MB`（参照の局所性テスト用に意図的に小さく設定）
- `schema.sql` が自動実行されてテーブルが作成される

### 2. データセットアップ（100万レコード投入）

```bash
go test -v -run TestSetup -timeout 10m
```

`TestSetupSerialData` と `TestSetupUUIDData` がそれぞれ100万件を1000件単位のバッチで INSERT する。

## ベンチマーク実行

### INSERT ベンチマーク

```bash
go test -bench BenchmarkInsert -benchmem -count 5
```

| ベンチマーク | 内容 |
|---|---|
| `BenchmarkInsertSerial` | `users_serial` へ INSERT（id は DB 自動採番） |
| `BenchmarkInsertUUID` | `users_uuid` へ INSERT（Go 側で UUID v7 生成） |

### 参照の局所性ベンチマーク

事前に `TestSetup` で100万件を投入してから実行する。

```bash
go test -bench BenchmarkSelect -benchmem -count 5
```

| ベンチマーク | 内容 |
|---|---|
| `BenchmarkSelectRecentSerial` | 最近の1000件からランダムに SELECT |
| `BenchmarkSelectRecentUUID` | 最近の1000件からランダムに SELECT |
| `BenchmarkSelectRandomSerial` | 全範囲からランダムに SELECT |
| `BenchmarkSelectRandomUUID` | 全範囲からランダムに SELECT |

## テーブルサイズの確認

psql で接続してサイズを確認できる。

```bash
docker compose exec postgres psql -U postgres -d bench
```

### テーブル本体・インデックスのサイズ

```sql
SELECT
    c.relname AS table_name,
    pg_size_pretty(pg_total_relation_size(c.oid)) AS total_size,
    pg_size_pretty(pg_relation_size(c.oid)) AS table_size,
    pg_size_pretty(pg_indexes_size(c.oid)) AS index_size,
    s.n_live_tup AS row_count
FROM pg_class c
JOIN pg_stat_user_tables s ON c.relname = s.relname
WHERE c.relname IN ('users_serial', 'users_uuid')
ORDER BY c.relname;
```

100万件投入時の参考値:

| table_name | total_size | table_size | index_size | row_count |
|---|---|---|---|---|
| users_serial | 102 MB | 80 MB | 21 MB | 1,000,000 |
| users_uuid | 119 MB | 89 MB | 30 MB | 1,000,000 |

UUID テーブルの方が全体で約17% 大きい（PK カラム: bigint 8 bytes vs UUID 16 bytes）。

## クリーンアップ

```bash
docker compose down -v
```
