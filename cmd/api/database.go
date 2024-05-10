package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func ConnectPgPoolConfigured(dsn *string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(*dsn)
	if err != nil {
		return nil, err
	}
	config.MinConns = 50
	config.MaxConns = 200
	config.MaxConnLifetime = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return pool, nil
}
