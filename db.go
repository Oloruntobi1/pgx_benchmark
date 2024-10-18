package pgx_benchmark

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	SqlxDB  *sqlx.DB
	PgxPool *pgxpool.Pool
}

func NewDB(dbURL string) (*DB, error) {
	sqlxDB, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	pgxPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	return &DB{
		SqlxDB:  sqlxDB,
		PgxPool: pgxPool,
	}, nil
}

func (db *DB) Close() {
	db.SqlxDB.Close()
	db.PgxPool.Close()
}
