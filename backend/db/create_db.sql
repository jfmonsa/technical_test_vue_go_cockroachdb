-- This SQL script creates the necessary tables for app, must be executed once.
-- Before running this script, ensure that the database is created and the connection is established.
-- The script is idempotent, meaning it can be run multiple times without causing errors.
CREATE TABLE IF NOT EXISTS stocks (
    ticker TEXT NOT NULL,
    company TEXT,
    brokerage TEXT,
    action TEXT,
    rating_from TEXT,
    rating_to TEXT,
    target_from FLOAT,
    target_to FLOAT,
    time TIMESTAMP,
    recommendation_score FLOAT,
    PRIMARY KEY (ticker, time)
);

-- This table stores the raw JSON data for items that failed in TRANSFORM or LOAD phases of ETL process.
CREATE TABLE IF NOT EXISTS failed_items (
    id SERIAL PRIMARY KEY,
    raw_json JSONB NOT NULL,
    error_message TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    failed_at_phase TEXT NOT NULL
);
