package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// centralize db connection, easy to replace with aws RDS, env vars
func NewPostgresPool() *pgxpool.Pool {
	dsn := "postgres://notif_user:notif_pass@localhost:5433/notification_db"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("db not reachable", err)
	}
	return pool
}
