package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/sahil_malakar/production_grade_golang_setup/internal/middleware"
)

// acting as a data structure to  append
// & store all the scanned listing
type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

// Converted closure factory patter to constructor factory pattern
// to fix dependency injection problem
// listing class
type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

// contructor of listing class
// --> go idioms name of contructor start with New
func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		// private readonly
		db:     db,
		logger: logger,
	}
}

// now to attach the http handler with the construction
// we use golang methods
// in general this fuctions are termed as method function in opps terminology
func (this ListingHandler) Get(w http.ResponseWriter, r *http.Request) {
	this.logger.Debug("getting listings")

	// Use the request-scoped context so the database query
	// is cancelled if the client disconnects or the request ends.
	ctx := r.Context()
	reqId := middleware.RequestIDFromContext(ctx)

	rows, err := this.db.QueryContext(
		ctx,
		`SELECT id, title, description, price, city, created_at
			 FROM listings
			 ORDER BY created_at DESC
			 LIMIT 50`)
	if err != nil {
		this.logger.Error("db.Query failed", "error", err, "request_id", reqId)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	this.logger.Debug("db.Query executed successfully", "request_id", reqId)
	// run at the last to close the connections
	defer rows.Close()

	listings := make([]listing, 0, 50)
	for rows.Next() {
		var l listing
		if err := rows.Scan(
			&l.ID,
			&l.Title,
			&l.Description,
			&l.Price,
			&l.City,
			&l.CreatedAt,
		); err != nil {
			this.logger.Error("rows.Scan failed", "error", err, "request_id", reqId)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}
	if err := rows.Err(); err != nil {
		this.logger.Error("rows.Err failed", "error", err, "request_id", reqId)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	this.logger.Info("listings fetched successfully", "count", len(listings), "request_id", reqId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Serializes the listings slice into JSON
	// and sends it to the client.
	// send error automatically, if caught.
	_ = json.NewEncoder(w).Encode(listings)
}

func (this ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	this.logger.Debug("deleting listing")

	ctx := r.Context()
	reqId := middleware.RequestIDFromContext(ctx)

	id := r.PathValue("id")

	this.logger.Debug("listing id extracted", "id", id, "request_id", reqId)
	_, err := this.db.ExecContext(
		ctx,
		`DELETE FROM listings
			WHERE id = $1`, id,
	)

	if err != nil {
		this.logger.Error("delete failed", "error", err, "id", id, "request_id", reqId)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	this.logger.Info("listing deleted successfully", "id", id, "request_id", reqId)

	w.WriteHeader(http.StatusNoContent)
}
