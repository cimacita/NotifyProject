package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(dsn string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Error pool creating: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Error ping: %v", err)
	}

	return pool
}
