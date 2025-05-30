-- This SQL script creates the necessary tables for app, must be excecuted once.
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
    PRIMARY KEY (ticker, time)
);
