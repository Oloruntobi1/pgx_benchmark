package pgx_benchmark

import (
	"context"
	"fmt"
	"testing"
)

const dbURL = "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"

func BenchmarkQuery(b *testing.B) {
	db, err := NewDB(dbURL)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	query := "SELECT 1"

	b.Run("sqlx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var result int
			err := db.SqlxDB.Get(&result, query)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("pgxpool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var result int
			err := db.PgxPool.QueryRow(context.Background(), query).Scan(&result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkConcurrentQueries(b *testing.B) {
	db, err := NewDB(dbURL)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	query := "SELECT pg_sleep(0.01)" // Simulate a query that takes 10ms

	for _, concurrency := range []int{1, 10, 50, 100} {
		b.Run(fmt.Sprintf("sqlx-concurrency-%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := db.SqlxDB.Exec(query)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})

		b.Run(fmt.Sprintf("pgxpool-concurrency-%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := db.PgxPool.Exec(context.Background(), query)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkMixedWorkload(b *testing.B) {
	db, err := NewDB(dbURL)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	// Prepare the database
	_, err = db.SqlxDB.Exec(`
		CREATE TABLE IF NOT EXISTS benchmark_users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		b.Fatal(err)
	}

	for _, lib := range []string{"sqlx", "pgxpool"} {
		b.Run(lib, func(b *testing.B) {
			b.SetParallelism(50)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					switch lib {
					case "sqlx":
						// Insert
						_, err := db.SqlxDB.Exec("INSERT INTO benchmark_users (name) VALUES ($1)", "John Doe")
						if err != nil {
							b.Fatal(err)
						}
						// Select
						var count int
						err = db.SqlxDB.Get(&count, "SELECT COUNT(*) FROM benchmark_users")
						if err != nil {
							b.Fatal(err)
						}
					case "pgxpool":
						// Insert
						_, err := db.PgxPool.Exec(context.Background(), "INSERT INTO benchmark_users (name) VALUES ($1)", "John Doe")
						if err != nil {
							b.Fatal(err)
						}
						// Select
						var count int
						err = db.PgxPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM benchmark_users").Scan(&count)
						if err != nil {
							b.Fatal(err)
						}
					}
				}
			})
		})
	}

	// Clean up
	_, err = db.SqlxDB.Exec("DROP TABLE benchmark_users")
	if err != nil {
		b.Fatal(err)
	}
}
