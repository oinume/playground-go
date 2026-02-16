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

	// Get the max ID to determine "recent" range
	var maxID int64
	err := pool.QueryRow(ctx, "SELECT COALESCE(MAX(id), 0) FROM users_serial").Scan(&maxID)
	if err != nil {
		b.Fatal(err)
	}
	if maxID < 1000 {
		b.Skip("Not enough data in users_serial. Run TestSetupSerialData first.")
	}

	// Recent 1000 IDs
	minRecent := maxID - 999

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetID := minRecent + rand.Int64N(1000)
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

	// Get the most recent 1000 UUIDs
	rows, err := pool.Query(ctx,
		"SELECT id FROM users_uuid ORDER BY id DESC LIMIT 1000")
	if err != nil {
		b.Fatal(err)
	}
	var recentIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			b.Fatal(err)
		}
		recentIDs = append(recentIDs, id)
	}
	rows.Close()
	if len(recentIDs) < 1000 {
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

	var maxID int64
	err := pool.QueryRow(ctx, "SELECT COALESCE(MAX(id), 0) FROM users_serial").Scan(&maxID)
	if err != nil {
		b.Fatal(err)
	}
	if maxID < 1000 {
		b.Skip("Not enough data in users_serial. Run TestSetupSerialData first.")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetID := rand.Int64N(maxID) + 1
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

	// Collect a sample of UUIDs from across the full range
	rows, err := pool.Query(ctx,
		"SELECT id FROM users_uuid ORDER BY random() LIMIT 10000")
	if err != nil {
		b.Fatal(err)
	}
	var allIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			b.Fatal(err)
		}
		allIDs = append(allIDs, id)
	}
	rows.Close()
	if len(allIDs) < 1000 {
		b.Skip("Not enough data in users_uuid. Run TestSetupUUIDData first.")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetID := allIDs[rand.IntN(len(allIDs))]
		var name, email string
		err := pool.QueryRow(ctx,
			"SELECT name, email FROM users_uuid WHERE id = $1", targetID,
		).Scan(&name, &email)
		if err != nil {
			b.Fatal(err)
		}
	}
}
