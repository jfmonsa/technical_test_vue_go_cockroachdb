// Package etl implements an ETL (Extract, Transform, Load) process for stock data.
// This script fetches stock recommendation data from an external API, transforms
// the data into the required format, and loads it into a CockroachDB database.
package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"
	"vue_go_cockroachdb/app"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5"
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
				continue
			}
			err = insertStockItem(ctx, conn, item)
			if err != nil {
				log.Println("Insert error:", err)
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
func transform(raw APIRawItem) (StockItem, error) {
	targetFrom, err := parseDollar(raw.TargetFrom)
	if err != nil {
		return StockItem{}, err
	}
	targetTo, err := parseDollar(raw.TargetTo)
	if err != nil {
		return StockItem{}, err
	}
	parsedTime, err := time.Parse(time.RFC3339Nano, raw.Time)
	if err != nil {
		return StockItem{}, err
	}

	return StockItem{
		Ticker:     raw.Ticker,
		Company:    raw.Company,
		Brokerage:  raw.Brokerage,
		Action:     raw.Action,
		RatingFrom: raw.RatingFrom,
		RatingTo:   raw.RatingTo,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Time:       parsedTime,
	}, nil
}

// parseDollar removes the dollar sign from a string and parses it as a float64.
func parseDollar(s string) (float64, error) {
	s = strings.ReplaceAll(s, "$", "")
	return strconv.ParseFloat(s, 64)
}

// insertStockItem inserts a StockItem into the stocks table.
// If a record with the same ticker and time already exists, it does nothing.
func insertStockItem(ctx context.Context, conn *pgx.Conn, item StockItem) error {
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
