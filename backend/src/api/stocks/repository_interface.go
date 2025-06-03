package stocks

import (
	"context"
	"vue_go_cockroachdb/src/models"
)

// interface
type StockRepository interface {
	GetStocks(ctx context.Context, search, sortBy, order string, page, limit int) ([]models.Stock, int, error)
	GetStockByTicker(ctx context.Context, ticker string) (*models.Stock, error)
	GetTopRecommendedStocks(ctx context.Context, limit int, minimunScore float64) ([]StockWithScore, error)
}
