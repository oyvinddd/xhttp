package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, cfg.PostgresURI)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		// cleanup
		db.Close()
		return nil, err
	}
	return db, nil
}

