// main.go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq" // cockroach driver

	"vue_go_cockroachdb/src/api/stocks"
	"vue_go_cockroachdb/src/app"
)

// main is the entry point of the application. It initializes the database connection,
// sets up the stock repository and HTTP handlers, configures the router with endpoints
// for retrieving stocks and stock details by ticker, and starts the HTTP server.
func main() {
	db := app.GetDBConnection(context.Background())

	repo := stocks.NewCockroachDBStockRepository(db)
	handler := &stocks.Handler{Repo: repo}

	r := chi.NewRouter()

	// CORS middleware to allow cross-origin requests
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// to test:
	// curl "http://localhost:8080/stocks?page=1&limit=5"
	r.Get("/stocks", handler.GetStocks)

	// to test:
	// curl "http://localhost:8080/stocks/AKBA"
	r.Get("/stocks/{ticker}", handler.GetStockByTicker)

	// to test:
	// curl "http://localhost:8080/recommendations?limit=5&minimun_score=7
	r.Get("/recommendations", handler.GetRecommendations)

	log.Println("🚀 Server listening ")
	http.ListenAndServe(":"+app.EnvVarsValues.Port, r)
}
