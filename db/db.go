package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oyvinddd/xhttp/env"
)

func Connect(ctx context.Context, dbURI string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return nil, err
	}
	// we try pinging the db in order to check the connection
	if err := db.Ping(ctx); err != nil {
		// do cleanup if we fail to connect
		db.Close()
		return nil, err
	}

	return db, nil
}

func ConnectUsingEnv(ctx context.Context) (*pgxpool.Pool, error) {
	// For this to work we need the following in our .env:
	// POSTGRES_UR, POSTGRES_MAX_CONNS, POSTGRES_MIN_CONNS
	dbURI := env.StringOrFatal("POSTGRES_URI")
	dbMaxConns := env.IntOrFatal("POSTGRES_MAX_CONNS")
	dbMinConns := env.IntOrFatal("POSTGRES_MIN_CONNS")

	config, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(dbMaxConns)
	config.MinConns = int32(dbMinConns)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
