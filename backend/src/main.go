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

	// to test:
	// curl "http://localhost:8080/stocks?page=1&limit=5"
	r.Get("/stocks", handler.GetStocks)

	// to test:
	// curl "http://localhost:8080/stocks/AKBA"
	r.Get("/stocks/{ticker}", handler.GetStockByTicker)

	// to test:
	// curl "http://localhost:8080/recommendations"
	r.Get("/recommendations", handler.GetRecommendations)

	log.Println("🚀 Server listening ")
	http.ListenAndServe(":"+app.EnvVarsValues.Port, r)
}
