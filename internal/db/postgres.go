package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dsn = "postgres://postgres:password@localhost:5432/postgres"
	}

	return pgxpool.New(context.Background(), dsn)
}
