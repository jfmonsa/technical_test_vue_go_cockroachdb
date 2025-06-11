package app

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// GetDBConnection establishes and returns a connection to the PostgreSQL database
// using the connection string specified in EnvVarsValues.DBURL. If the connection
// cannot be established, the function logs a fatal error and terminates the application.
//
// Parameters:
//   - ctx: context.Context for managing request-scoped values, cancellation, and timeouts.
func GetDBConnection(ctx context.Context) *sql.DB {
	conn, err := sql.Open("postgres", EnvVarsValues.DB_URL)
	if err != nil {
		log.Fatal("DB Connection Error:", err)
	}
	return conn
}
