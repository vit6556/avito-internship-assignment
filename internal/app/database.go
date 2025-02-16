package app

import (
	"context"
	"log"
	"net"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vit6556/avito-internship-assignment/internal/config"
)

func InitDatabase() *pgxpool.Pool {
	ctx := context.Background()
	cfg := config.LoadDatabaseConfig()

	connString := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   cfg.Name,
	}

	dbPool, err := pgxpool.New(ctx, connString.String())
	if err != nil {
		log.Fatal("failed to create db pool")
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping db: %s", err.Error())
	}

	return dbPool
}
