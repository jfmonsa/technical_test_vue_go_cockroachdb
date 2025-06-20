// Package etl implements an ETL (Extract, Transform, Load) process for stock data.
// This script fetches stock recommendation data from an external API, transforms
// the data into the required format, and loads it into a CockroachDB database.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"vue_go_cockroachdb/src/app"
	"vue_go_cockroachdb/src/models"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5"
)

const (
	failedPhaseTransform = "TRANSFORM"
	failedPhaseLoad      = "LOAD" // insert in ETL terminology
)

// writeLogs initializes the logging system to write logs to a file named "etl.log".
// If the file does not exist, it will be created. If it exists, logs will be appended.
// It returns the file handle for the log file.
func writeLogs() *os.File {
	dateTimeStr := time.Now().Format(time.RFC3339)
	logDir := "logs"
	// Ensure the logs directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal("Could not create logs directory:", err)
	}
	logFileName := fmt.Sprintf("%s/etl-%s.log", logDir, dateTimeStr)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Could not open log file:", err)
	}

	log.SetOutput(logFile)
	return logFile
}

// main is the entry point for the ETL process.
// It connects to the database, fetches paginated stock data from an API,
// transforms each item, and inserts it into the database.
func main() {
	logFile := writeLogs()
	defer logFile.Close()

	ctx := context.Background()

	// Conectar a CockroachDB
	conn, err := pgx.Connect(ctx, app.EnvVarsValues.DB_URL)
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
				if err := insertFailedItem(ctx, conn, raw, err, failedPhaseTransform); err != nil {
					log.Println("Failed to insert failed item:", err)
				}
				continue
			}
			err = insertStockItem(ctx, conn, item)
			if err != nil {
				log.Println("Insert error:", err)
				if err := insertFailedItem(ctx, conn, raw, err, failedPhaseLoad); err != nil {
					log.Println("Failed to insert good item:", err)
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
func transform(raw APIRawItem) (models.StockWithScore, error) {
	if raw.Ticker == "" {
		return models.StockWithScore{}, fmt.Errorf("ticker is required but was empty")
	}
	if raw.Time == "" {
		return models.StockWithScore{}, fmt.Errorf("time is required but was empty for ticker '%s'", raw.Ticker)
	}

	// NOTE: there are registers that have an empty rating_from or rating_to the decision is to ignore them
	// because they are could be considered as "not rated" or "no recommendation" and bias the results.
	if !isValidRating(raw.RatingFrom) {
		return models.StockWithScore{}, fmt.Errorf("invalid rating_from value '%s' for ticker '%s'", raw.RatingFrom, raw.Ticker)
	}
	if !isValidRating(raw.RatingTo) {
		return models.StockWithScore{}, fmt.Errorf("invalid rating_to value '%s' for ticker '%s'", raw.RatingTo, raw.Ticker)
	}

	targetFrom, err := parseDollar(raw.TargetFrom)
	if err != nil {
		return models.StockWithScore{}, fmt.Errorf("invalid target_from value '%s' for ticker '%s': %v", raw.TargetFrom, raw.Ticker, err)
	}
	targetTo, err := parseDollar(raw.TargetTo)
	if err != nil {
		return models.StockWithScore{}, fmt.Errorf("invalid target_to value '%s' for ticker '%s': %v", raw.TargetTo, raw.Ticker, err)
	}

	stockStruct := models.Stock{
		Ticker:     raw.Ticker,
		Company:    raw.Company,
		Brokerage:  raw.Brokerage,
		Action:     raw.Action,
		RatingFrom: raw.RatingFrom,
		RatingTo:   raw.RatingTo,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Time:       raw.Time,
	}

	score := CalculateStockScore(stockStruct)

	// assign the score to a new struct that includes the stock and the score
	var stockStructWithScore models.StockWithScore
	stockStructWithScore.Stock = stockStruct
	stockStructWithScore.RecommendationScore = score
	return stockStructWithScore, nil
}

// parseDollar removes the dollar sign from a string and parses it as a float64.
func parseDollar(s string) (float64, error) {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "") // Delete commas for thousands separators
	return strconv.ParseFloat(s, 64)
}

// insertStockItem inserts a StockItem into the stocks table.
// If a record with the same ticker and time already exists, it does nothing.
func insertStockItem(ctx context.Context, conn *pgx.Conn, item models.StockWithScore) error {
	_, err := conn.Exec(ctx, `
		INSERT INTO stocks (
			ticker, company, brokerage, action, rating_from, rating_to,
			target_from, target_to, time, recommendation_score
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9, $10)
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
		item.RecommendationScore,
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
