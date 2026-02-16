package postgres_serial_vs_uuid

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:pass@localhost:5433/bench?sslmode=disable"
	}
	fmt.Printf("Starting benchmarks: dsn=%v\n", dsn)

	var err error
	pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	os.Exit(m.Run())
}

// --- INSERT benchmarks ---

func BenchmarkInsertSerial(b *testing.B) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE users_serial")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pool.Exec(ctx,
			"INSERT INTO users_serial (name, email) VALUES ($1, $2)",
			fmt.Sprintf("user-%d", i),
			fmt.Sprintf("user-%d@example.com", i),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInsertUUID(b *testing.B) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE users_uuid")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, err := uuid.NewV7()
		if err != nil {
			b.Fatal(err)
		}
		_, err = pool.Exec(ctx,
			"INSERT INTO users_uuid (id, name, email) VALUES ($1, $2, $3)",
			id,
			fmt.Sprintf("user-%d", i),
			fmt.Sprintf("user-%d@example.com", i),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- Setup functions for locality benchmarks ---

const setupRows = 1_000_000
const batchSize = 1000
const recentSampleSize = 1000
const randomSampleSize = 10000

func TestSetupSerialData(t *testing.T) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE users_serial RESTART IDENTITY")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < setupRows; i += batchSize {
		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		for j := 0; j < batchSize && i+j < setupRows; j++ {
			n := i + j
			_, err := tx.Exec(ctx,
				"INSERT INTO users_serial (name, email) VALUES ($1, $2)",
				fmt.Sprintf("user-%d", n),
				fmt.Sprintf("user-%d@example.com", n),
			)
			if err != nil {
				tx.Rollback(ctx)
				t.Fatal(err)
			}
		}
		if err := tx.Commit(ctx); err != nil {
			t.Fatal(err)
		}
		if (i+batchSize)%100_000 == 0 {
			t.Logf("Inserted %d rows into users_serial", i+batchSize)
		}
	}
	t.Logf("Setup complete: %d rows in users_serial", setupRows)
}

func TestSetupUUIDData(t *testing.T) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE users_uuid RESTART IDENTITY")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < setupRows; i += batchSize {
		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		for j := 0; j < batchSize && i+j < setupRows; j++ {
			n := i + j
			id, err := uuid.NewV7()
			if err != nil {
				tx.Rollback(ctx)
				t.Fatal(err)
			}
			_, err = tx.Exec(ctx,
				"INSERT INTO users_uuid (id, name, email) VALUES ($1, $2, $3)",
				id,
				fmt.Sprintf("user-%d", n),
				fmt.Sprintf("user-%d@example.com", n),
			)
			if err != nil {
				tx.Rollback(ctx)
				t.Fatal(err)
			}
		}
		if err := tx.Commit(ctx); err != nil {
			t.Fatal(err)
		}
		if (i+batchSize)%100_000 == 0 {
			t.Logf("Inserted %d rows into users_uuid", i+batchSize)
		}
	}
	t.Logf("Setup complete: %d rows in users_uuid", setupRows)
}

// --- SELECT benchmarks (locality) ---

func BenchmarkSelectRecentSerial(b *testing.B) {
	ctx := context.Background()

	recentIDs, err := loadRecentSerialIDs(ctx, recentSampleSize)
	if err != nil {
		b.Fatal(err)
	}
	if len(recentIDs) < recentSampleSize {
		b.Skip("Not enough data in users_serial. Run TestSetupSerialData first.")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetID := recentIDs[rand.IntN(len(recentIDs))]
		var name, email string
		err := pool.QueryRow(ctx,
			"SELECT name, email FROM users_serial WHERE id = $1", targetID,
		).Scan(&name, &email)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSelectRecentUUID(b *testing.B) {
	ctx := context.Background()

	recentIDs, err := loadRecentUUIDIDs(ctx, recentSampleSize)
	if err != nil {
		b.Fatal(err)
	}
	if len(recentIDs) < recentSampleSize {
		b.Skip("Not enough data in users_uuid. Run TestSetupUUIDData first.")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetID := recentIDs[rand.IntN(len(recentIDs))]
		var name, email string
		err := pool.QueryRow(ctx,
			"SELECT name, email FROM users_uuid WHERE id = $1", targetID,
		).Scan(&name, &email)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSelectRandomSerial(b *testing.B) {
	ctx := context.Background()

	randomIDs, err := loadRandomSerialIDs(ctx, randomSampleSize)
	if err != nil {
		b.Fatal(err)
	}
	if len(randomIDs) < recentSampleSize {
		b.Skip("Not enough data in users_serial. Run TestSetupSerialData first.")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetID := randomIDs[rand.IntN(len(randomIDs))]
		var name, email string
		err := pool.QueryRow(ctx,
			"SELECT name, email FROM users_serial WHERE id = $1", targetID,
		).Scan(&name, &email)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSelectRandomUUID(b *testing.B) {
	ctx := context.Background()

	randomIDs, err := loadRandomUUIDIDs(ctx, randomSampleSize)
	if err != nil {
		b.Fatal(err)
	}
	if len(randomIDs) < recentSampleSize {
		b.Skip("Not enough data in users_uuid. Run TestSetupUUIDData first.")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetID := randomIDs[rand.IntN(len(randomIDs))]
		var name, email string
		err := pool.QueryRow(ctx,
			"SELECT name, email FROM users_uuid WHERE id = $1", targetID,
		).Scan(&name, &email)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func loadRecentSerialIDs(ctx context.Context, limit int) ([]int64, error) {
	rows, err := pool.Query(ctx,
		"SELECT id FROM users_serial ORDER BY id DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0, limit)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func loadRecentUUIDIDs(ctx context.Context, limit int) ([]uuid.UUID, error) {
	rows, err := pool.Query(ctx,
		"SELECT id FROM users_uuid ORDER BY id DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]uuid.UUID, 0, limit)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func loadRandomSerialIDs(ctx context.Context, limit int) ([]int64, error) {
	ids := make([]int64, 0, limit)
	samplePercents := []string{"1.0", "2.0", "5.0", "10.0"}
	for _, percent := range samplePercents {
		query := fmt.Sprintf(
			"SELECT id FROM users_serial TABLESAMPLE SYSTEM (%s) LIMIT %d",
			percent, limit,
		)
		rows, err := pool.Query(ctx, query)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return nil, err
			}
			ids = append(ids, id)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, err
		}
		rows.Close()
		if len(ids) >= limit {
			break
		}
	}
	return ids, nil
}

func loadRandomUUIDIDs(ctx context.Context, limit int) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, limit)
	samplePercents := []string{"1.0", "2.0", "5.0", "10.0"}
	for _, percent := range samplePercents {
		query := fmt.Sprintf(
			"SELECT id FROM users_uuid TABLESAMPLE SYSTEM (%s) LIMIT %d",
			percent, limit,
		)
		rows, err := pool.Query(ctx, query)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return nil, err
			}
			ids = append(ids, id)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, err
		}
		rows.Close()
		if len(ids) >= limit {
			break
		}
	}
	return ids, nil
}
