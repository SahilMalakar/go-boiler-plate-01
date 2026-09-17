package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sahil_malakar/production_grade_golang_setup/internal/config"
)

func DbConnect(databaseURL string, cfg config.DBConfig) (*sql.DB, error) {

	// Creates the database handle using the PostgreSQL connection URL.
	// The handle manages the connection pool used by the application.
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Limits the maximum number of connections that can be opened
	// simultaneously by the connection pool.
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	// Keeps the configured number of unused connections available
	// so they can be reused instead of creating new connections.
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	// Closes connections that have existed longer than the configured lifetime.
	// New connections will be created when the pool needs them.
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Creates a context that limits how long the initial database
	// connectivity check is allowed to run.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.PingTimeout,
	)
	defer cancel()

	// Verifies that the application can actually communicate with PostgreSQL.
	// The context prevents this check from waiting indefinitely.
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.ping: %w", err)
	}

	return db, nil
}
