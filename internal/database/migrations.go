package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// Migrate erstellt die Tabelle `offers`.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {

	log.Print("Running database migrations...")

	query := `
	CREATE TABLE IF NOT EXISTS offers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
