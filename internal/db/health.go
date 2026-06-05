package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Status represents the database connection state.
type Status string

const (
	StatusOK          Status = "ok"
	StatusDisabled    Status = "disabled"
	StatusUnreachable Status = "unreachable"
)

// HealthResult holds the result of a database health check.
type HealthResult struct {
	Status   Status `json:"status"`
	Database string `json:"database"`
}

// Ping checks whether the database pool is reachable.
// pool == nil means DB is not configured.
func Ping(pool *pgxpool.Pool) HealthResult {
	if pool == nil {
		return HealthResult{Status: StatusDisabled, Database: "not_configured"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return HealthResult{Status: StatusUnreachable, Database: "unreachable"}
	}

	return HealthResult{Status: StatusOK, Database: "connected"}
}
