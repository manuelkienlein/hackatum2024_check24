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
    region_id INTEGER NOT NULL, -- Region ID
    time_range_start BIGINT NOT NULL, -- Start time of the range (ms since UNIX epoch)
    time_range_end BIGINT NOT NULL, -- End time of the range (ms since UNIX epoch)
    number_days INTEGER NOT NULL, -- Number of full days available
    sort_order VARCHAR(20) NOT NULL CHECK (sort_order IN ('price-asc', 'price-desc')), -- Sort order (price ascending or descending)
    page INTEGER NOT NULL, -- Pagination page number
    page_size INTEGER NOT NULL, -- Number of offers per page
    price_range_width INTEGER NOT NULL, -- Price range width in cents
    min_free_kilometer_width INTEGER NOT NULL, -- Minimum free kilometer range width in km
    min_number_seats INTEGER, -- Minimum number of seats in the car
    min_price NUMERIC(10, 2), -- Minimum price in cents
    max_price NUMERIC(10, 2), -- Maximum price in cents
    car_type VARCHAR(20) CHECK (car_type IN ('small', 'sports', 'luxury', 'family')), -- Type of the car
    only_vollkasko BOOLEAN NOT NULL, -- Whether only offers with vollkasko are included
    min_free_kilometer INTEGER -- Minimum free kilometers included
	);
	
	-- Indexing suggestions for faster query performance
	CREATE INDEX idx_offers_region_time ON offers (region_id, time_range_start, time_range_end);
	CREATE INDEX idx_offers_price ON offers (min_price, max_price);
	CREATE INDEX idx_offers_type ON offers (car_type);`

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
