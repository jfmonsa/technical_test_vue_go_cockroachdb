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

func (r *CockroachDBStockRepository) GetStocks(ctx context.Context, search, sortBy, order string, page, limit int) ([]models.Stock, error) {
	offset := (page - 1) * limit
	baseQuery := `
        SELECT ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time
        FROM stocks
    `

	var filters []string
	var args []interface{}
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
		return nil, err
	}
	defer rows.Close()
	var stocks []models.Stock
	for rows.Next() {
		var s models.Stock
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
		)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
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
