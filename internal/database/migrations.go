package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// Create offers table
	query := `
	CREATE TABLE offers (
    id SERIAL PRIMARY KEY, -- Unique identifier for each offer
    data VARCHAR(256) NOT NULL, -- additional data of the offer
    most_specific_region_id INTEGER NOT NULL, -- Region ID
    start_date BIGINT NOT NULL, -- Start time of the range (ms since UNIX epoch)
    end_date BIGINT NOT NULL, -- End time of the range (ms since UNIX epoch)
    number_Seats INTEGER NOT NULL, -- Number of seats in the car
    price INTEGER NOT NULL, -- Price in cents
    number_days INTEGER NOT NULL, -- Number of full days available
    car_type VARCHAR(20), -- Type of the car
    only_vollkasko BOOLEAN NOT NULL, -- Whether only offers with vollkasko are included
    free_kilometers INTEGER -- free kilometers included
	);`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	// Create static_region_data table with parent_id
	query = `
	CREATE TABLE IF NOT EXISTS static_region_data (
		id INT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		parent_id INT,
		FOREIGN KEY (parent_id) REFERENCES static_region_data(id)
	)`

	_, err = pool.Exec(ctx, query)
	if err != nil {
		return err
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

	byteValue, err := ioutil.ReadAll(file)
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

// DropTables drops the offers and static_region_data tables/fix
func DropTables(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "DROP TABLE IF EXISTS offers")
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, "DROP TABLE IF EXISTS static_region_data")
	if err != nil {
		return err
	}

	return nil
}
