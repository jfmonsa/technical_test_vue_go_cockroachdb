// internal/stocks/handler.go
package stocks

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	Repo StockRepository
}

func (h *Handler) GetStocks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	search := q.Get("search")
	sortBy := q.Get("sort_by")
	order := q.Get("order")
	if order == "" {
		order = "desc"
	}

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 {
		limit = 10
	}

	stocks, total, err := h.Repo.GetStocks(r.Context(), search, sortBy, order, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	resp := map[string]interface{}{
		"items":      stocks,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetStockByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")

	stock, err := h.Repo.GetStockByTicker(r.Context(), ticker)
	if err != nil {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 5
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	minimunScore := 0.0
	if s := r.URL.Query().Get("minimun_score"); s != "" {
		if parsed, err := strconv.ParseFloat(s, 64); err == nil && parsed >= 0 {
			minimunScore = parsed
		}
	}

	stocks, err := h.Repo.GetTopRecommendedStocks(ctx, limit, minimunScore)
	if err != nil {
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}
