// Package etl implements an ETL (Extract, Transform, Load) process for stock data.
// This script fetches stock recommendation data from an external API, transforms
// the data into the required format, and loads it into a CockroachDB database.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"vue_go_cockroachdb/src/app"
	"vue_go_cockroachdb/src/models"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5"
)

const (
	phailedPhaseTransform = "TRANSFORM"
	phailedPhaseLoad      = "LOAD" // insert in ETL terminology
)

// main is the entry point for the ETL process.
// It connects to the database, fetches paginated stock data from an API,
// transforms each item, and inserts it into the database.
func main() {
	ctx := context.Background()

	// Conectar a CockroachDB
	conn, err := pgx.Connect(ctx, app.EnvVarsValues.DBURL)
	if err != nil {
		log.Fatal("DB Connection Error:", err)
	}
	defer conn.Close(ctx)

	client := resty.New()

	nextPage := ""
	for {
		resp, err := client.R().
			SetHeader("Authorization", app.EnvVarsValues.AuthToken).
			SetHeader("Content-Type", "application/json").
			SetQueryParam("next_page", nextPage).
			SetResult(&APIResponse{}).
			Get(app.EnvVarsValues.ApiURL)

		if err != nil {
			log.Fatal("API request failed:", err)
		}

		apiResp := resp.Result().(*APIResponse)

		for _, raw := range apiResp.Items {
			item, err := transform(raw)
			if err != nil {
				log.Println("Skipping item due to error:", err)
				if err := insertFailedItem(ctx, conn, raw, err, phailedPhaseTransform); err != nil {
					log.Println("Failed to insert failed item:", err)
				}
				continue
			}
			err = insertStockItem(ctx, conn, item)
			if err != nil {
				log.Println("Insert error:", err)
				if err := insertFailedItem(ctx, conn, raw, err, phailedPhaseLoad); err != nil {
					log.Println("Failed to insert failed item:", err)
				}
			}
		}

		if apiResp.NextPage == "" {
			break
		}
		nextPage = apiResp.NextPage
	}
}

// transform converts a raw API item into a StockItem struct,
// parsing dollar values and timestamps as needed.
func transform(raw APIRawItem) (models.Stock, error) {
	if !isValidRating(raw.RatingFrom) {
		return models.Stock{}, fmt.Errorf("invalid rating_from value '%s' for ticker '%s'", raw.RatingFrom, raw.Ticker)
	}
	if !isValidRating(raw.RatingTo) {
		return models.Stock{}, fmt.Errorf("invalid rating_to value '%s' for ticker '%s'", raw.RatingTo, raw.Ticker)
	}

	targetFrom, err := parseDollar(raw.TargetFrom)
	if err != nil {
		return models.Stock{}, fmt.Errorf("invalid target_from value '%s' for ticker '%s': %v", raw.TargetFrom, raw.Ticker, err)
	}
	targetTo, err := parseDollar(raw.TargetTo)
	if err != nil {
		return models.Stock{}, fmt.Errorf("invalid target_to value '%s' for ticker '%s': %v", raw.TargetTo, raw.Ticker, err)
	}

	return models.Stock{
		Ticker:     raw.Ticker,
		Company:    raw.Company,
		Brokerage:  raw.Brokerage,
		Action:     raw.Action,
		RatingFrom: raw.RatingFrom,
		RatingTo:   raw.RatingTo,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Time:       raw.Time,
	}, nil
}

// parseDollar removes the dollar sign from a string and parses it as a float64.
func parseDollar(s string) (float64, error) {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "") // Delete commas for thousands separators
	return strconv.ParseFloat(s, 64)
}

// insertStockItem inserts a StockItem into the stocks table.
// If a record with the same ticker and time already exists, it does nothing.
func insertStockItem(ctx context.Context, conn *pgx.Conn, item models.Stock) error {
	_, err := conn.Exec(ctx, `
		INSERT INTO stocks (
			ticker, company, brokerage, action, rating_from, rating_to,
			target_from, target_to, time
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (ticker, time) DO NOTHING
	`,
		item.Ticker,
		item.Company,
		item.Brokerage,
		item.Action,
		item.RatingFrom,
		item.RatingTo,
		item.TargetFrom,
		item.TargetTo,
		item.Time,
	)
	return err
}

// insertFailedItem inserts a the raw json of the failed item into the "failed_items" table in the db
// failed_at_phase indicates the phase of the ETL process where the failure occurred, can be "transform" or "insert".
func insertFailedItem(ctx context.Context, conn *pgx.Conn, raw APIRawItem, parseErr error, failed_at_phase string) error {
	rawJSON, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, `
        INSERT INTO failed_items (raw_json, error_message, failed_at_phase)
        VALUES ($1, $2,$3)
    `, string(rawJSON), parseErr.Error(), failed_at_phase)
	return err
}
