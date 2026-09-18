package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/sahil_malakar/production_grade_golang_setup/internal/config"
	"github.com/sahil_malakar/production_grade_golang_setup/internal/db"
	"github.com/sahil_malakar/production_grade_golang_setup/internal/handlers"
	"github.com/sahil_malakar/production_grade_golang_setup/internal/middleware"
)

func main() {
	// Loads application configuration from environment variables.
	cfg := config.MustLoad()

	// logger configurations
	logFile, err := os.OpenFile(
		"app.jsonl",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	writer := io.MultiWriter(
		os.Stdout,
		// logFile,
	)

	logHandler := slog.NewJSONHandler(
		writer,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		},
	)

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	logger.Debug("application configuration loaded")

	// Creates the PostgreSQL connection pool.
	db, err := db.DbConnect(cfg.DatabaseURL, cfg.DB)
	if err != nil {
		logger.Error("Database connection failed", "error", err)
	}
	logger.Info("Database connection established")
	// Closes the database pool when the application exits.
	defer db.Close()

	logger.Info("App server is running..")

	// initialized validator
	validate := validator.New()

	// Creates a router that receives incoming requests
	// and forwards them to the handler registered for the matching route.
	mux := http.NewServeMux()

	listingHandler := handlers.NewListingHandler(
		db,
		logger,
		validate,
	)

	// Registers the GET /health route and executes this function
	// whenever a request is made to that endpoint.
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /listings", listingHandler.Get)
	mux.HandleFunc("POST /listings", listingHandler.Create)
	mux.HandleFunc("DELETE /listings/{id}", listingHandler.Delete)

	// Creates the HTTP server and configures how it accepts and handles requests.
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      middleware.RequestId(mux), // wrapping middleware in the router handler
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}
	logger.Info("Server is listening on http://localhost:" + cfg.Port)

	// Starts the server and keeps the application running while it accepts requests.
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Server failed:", "error", err)
	}
}
