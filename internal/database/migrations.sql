-- Create offers table
CREATE TABLE IF NOT EXISTS offers (
    id UUID PRIMARY KEY, -- Unique identifier for each offer
    data VARCHAR(500) NOT NULL, -- additional data of the offer
    most_specific_region_id INTEGER NOT NULL, -- Region ID
    start_date BIGINT NOT NULL, -- Start time of the range (ms since UNIX epoch)
    end_date BIGINT NOT NULL, -- End time of the range (ms since UNIX epoch)
    number_Seats INTEGER NOT NULL, -- Number of seats in the car
    price INTEGER NOT NULL, -- Price in cents
    car_type VARCHAR(20), -- Type of the car
    only_vollkasko BOOLEAN NOT NULL, -- Whether only offers with vollkasko are included
    free_kilometers INTEGER -- free kilometers included
);

-- Create static_region_data table with parent_id
CREATE TABLE IF NOT EXISTS static_region_data (
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INT,
    FOREIGN KEY (parent_id) REFERENCES static_region_data(id)
);