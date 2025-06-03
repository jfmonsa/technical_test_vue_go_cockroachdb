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

type StockWithScore struct {
	models.Stock
	Score float64
}

func (r *CockroachDBStockRepository) GetTopRecommendedStocks(ctx context.Context, limit int, minimunScore float64) ([]StockWithScore, error) {
	query := `
        SELECT ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time
        FROM stocks
        ORDER BY time DESC
        `

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recommendations []StockWithScore

	for rows.Next() {
		var s models.Stock
		if err := rows.Scan(&s.Ticker, &s.Company, &s.Brokerage, &s.Action, &s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.Time); err != nil {
			continue
		}

		score := CalculateStockScore(s)

		// Format score for display in the company name
		// s.Company = fmt.Sprintf("%s (score: %.2f)", s.Company, score)

		// Umbral mínimo (opcional)
		if score >= minimunScore {
			recommendations = append(recommendations, StockWithScore{
				Stock: s,
				Score: score,
			})
		}
	}

	// Limitar resultados
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}
