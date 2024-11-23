package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Region represents the structure of the regions in the JSON file
type Region struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Subregions []Region `json:"subregions"`
}

// Migrate creates the tables and fills the static_region_data table
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {

	// Load SQL migrations from file
	migrationSQL, err := os.ReadFile("internal/database/migrations.sql")
	if err != nil {
		return fmt.Errorf("failed to read migrations.sql: %v", err)
	}

	// Execute the migration script
	_, err = pool.Exec(ctx, string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to execute migrations: %v", err)
	}

	// Read and parse regions.json
	file, err := os.Open("internal/database/regions.json")
	if err != nil {
		return fmt.Errorf("failed to open regions.json: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close regions.json: %v", err)
		}
	}(file)

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read regions.json: %v", err)
	}

	var rootRegion Region
	if err := json.Unmarshal(byteValue, &rootRegion); err != nil {
		return fmt.Errorf("failed to unmarshal regions.json: %v", err)
	}

	// Insert the root region and its subregions
	if err := insertRegion(ctx, pool, rootRegion, nil); err != nil {
		return fmt.Errorf("failed to insert root region: %v", err)
	}

	return nil
}

// insertRegion inserts a region and its subregions into the static_region_data table
func insertRegion(ctx context.Context, pool *pgxpool.Pool, region Region, parentID *int) error {
	_, err := pool.Exec(ctx, "INSERT INTO static_region_data (id, name, parent_id) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING", region.ID, region.Name, parentID)
	if err != nil {
		return err
	}

	for _, subregion := range region.Subregions {
		if err := insertRegion(ctx, pool, subregion, &region.ID); err != nil {
			return err
		}
	}

	return nil
}

func DropTables(ctx context.Context, pool *pgxpool.Pool) error {
	// List of tables to be dropped
	tables := []string{
		"offers",
		"static_region_data",
	}

	// Loop through each table and drop it
	for _, table := range tables {
		_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS "+table)
		if err != nil {
			return err
		}
	}

	return nil
}
