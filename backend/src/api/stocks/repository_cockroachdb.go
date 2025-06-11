package stocks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"vue_go_cockroachdb/src/models"
)

type CockroachDBStockRepository struct {
	DB *sql.DB
}

func NewCockroachDBStockRepository(db *sql.DB) *CockroachDBStockRepository {
	return &CockroachDBStockRepository{DB: db}
}

func (r *CockroachDBStockRepository) GetStocks(ctx context.Context, search, sortBy, order string, page, limit int) ([]models.Stock, int, error) {
	offset := (page - 1) * limit
	baseQuery := `
        SELECT ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time,
               COUNT(*) OVER() as total_count
        FROM stocks
    `

	var filters []string
	var args []any
	argIndex := 1

	if search != "" {
		filters = append(filters, fmt.Sprintf("(LOWER(ticker) LIKE LOWER($%d) OR LOWER(company) LIKE LOWER($%d))", argIndex, argIndex+1))
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	if len(filters) > 0 {
		baseQuery += " WHERE " + strings.Join(filters, " AND ")
	}

	if sortBy != "" {
		baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, strings.ToUpper(order))
	} else {
		baseQuery += " ORDER BY time DESC"
	}

	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.DB.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var stocks []models.Stock
	total := 0
	for rows.Next() {
		var s models.Stock
		var rowTotal int
		err := rows.Scan(
			&s.Ticker,
			&s.Company,
			&s.Brokerage,
			&s.Action,
			&s.RatingFrom,
			&s.RatingTo,
			&s.TargetFrom,
			&s.TargetTo,
			&s.Time,
			&rowTotal,
		)
		if err != nil {
			return nil, 0, err
		}
		stocks = append(stocks, s)
		total = rowTotal
	}
	return stocks, total, nil
}

func (r *CockroachDBStockRepository) GetStockByTicker(ctx context.Context, ticker string) (*models.Stock, error) {
	query := `
        SELECT ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time
        FROM stocks WHERE ticker = $1
    `
	row := r.DB.QueryRowContext(ctx, query, ticker)
	var s models.Stock
	err := row.Scan(
		&s.Ticker,
		&s.Company,
		&s.Brokerage,
		&s.Action,
		&s.RatingFrom,
		&s.RatingTo,
		&s.TargetFrom,
		&s.TargetTo,
		&s.Time,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *CockroachDBStockRepository) GetTopRecommendedStocks(ctx context.Context, page, limit int, minimumScore float64) ([]models.StockWithScore, error) {
	offset := (page - 1) * limit

	// esto porque ya todo esta calculado en la bd por tanto no hace falta calcularlo de nuevo
	query := `
        SELECT ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time, recommendation_score
        FROM stocks
		WHERE recommendation_score >= $1
        ORDER BY  recommendation_score time DESC
		LIMIT $2 OFFSET $3
        `

	rows, err := r.DB.QueryContext(ctx, query, minimumScore, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recommendations []models.StockWithScore

	for rows.Next() {
		var s models.StockWithScore
		if err := rows.Scan(&s.Ticker, &s.Company, &s.Brokerage, &s.Action, &s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.Time); err != nil {
			continue
		}

		// Format score for display in the company name
		s.Company = fmt.Sprintf("%s (score: %.2f)", s.Company, s.RecommendationScore)

		recommendations = append(recommendations, models.StockWithScore{
			Stock:               s.Stock,
			RecommendationScore: s.RecommendationScore,
		})

	}

	return recommendations, nil
}
