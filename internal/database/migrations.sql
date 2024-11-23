-- Create offers table
CREATE TABLE IF NOT EXISTS offers (
    id SERIAL PRIMARY KEY, -- Unique identifier for each offer
    data VARCHAR(256) NOT NULL, -- additional data of the offer
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
    car_type VARCHAR(20), -- Type of the car
    only_vollkasko BOOLEAN NOT NULL, -- Whether only offers with vollkasko are included
    min_free_kilometer INTEGER -- Minimum free kilometers included
);

-- Indexing suggestions for faster query performance
CREATE INDEX idx_offers_region_time ON offers (region_id, time_range_start, time_range_end);
CREATE INDEX idx_offers_price ON offers (min_price, max_price);
CREATE INDEX idx_offers_type ON offers (car_type);

-- Create static_region_data table with parent_id
CREATE TABLE IF NOT EXISTS static_region_data (
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INT,
    FOREIGN KEY (parent_id) REFERENCES static_region_data(id)
);