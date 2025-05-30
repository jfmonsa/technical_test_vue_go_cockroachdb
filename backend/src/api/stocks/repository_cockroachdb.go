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

func (r *CockroachDBStockRepository) GetTopRecommendedStocks(ctx context.Context, limit int) ([]models.Stock, error) {
	query := `
        SELECT ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time
        FROM stocks
        ORDER BY time DESC
        LIMIT 100
    `
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recommendations []models.Stock

	for rows.Next() {
		var s models.Stock
		if err := rows.Scan(&s.Ticker, &s.Company, &s.Brokerage, &s.Action, &s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.Time); err != nil {
			continue
		}

		// Simple scoring system
		score := 0.0

		// 1. Potencial de ganancia
		if s.TargetFrom > 0 {
			potential := ((s.TargetTo - s.TargetFrom) / s.TargetFrom) * 100
			if potential > 0 {
				score += potential / 5 // Normalizar
			}
		}

		// 2. Acción recomendada
		action := strings.ToLower(s.Action)
		if strings.Contains(action, "buy") || strings.Contains(action, "outperform") {
			score += 2
		}

		// 3. Mejora de calificación
		if s.RatingFrom != "" && s.RatingTo != "" && s.RatingFrom != s.RatingTo {
			score += 1
		}

		// 4. Reciente (deberías parsear fecha real)
		// Aquí deberías convertir s.Time a time.Time y comparar si es reciente
		// Solo para ejemplo:
		score += 1

		// Guardamos el score temporalmente como parte del ticker para debug
		s.Company = fmt.Sprintf("%s (score: %.2f)", s.Company, score)

		// Umbral mínimo (opcional)
		if score >= 3.0 {
			recommendations = append(recommendations, s)
		}
	}

	// Limitar resultados
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}
