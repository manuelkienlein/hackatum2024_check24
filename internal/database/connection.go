package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// ConnectDB stellt die Verbindung zu PostgreSQL mit pgxpool her.
func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:pg%23pass123@localhost:5432/template1?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return pool, nil
}
