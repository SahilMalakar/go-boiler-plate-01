package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
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
	db *sql.DB
}

// contructor of listing class
// --> go idioms name of contructor start with New
func NewListingHandler(db *sql.DB) *ListingHandler {
	return &ListingHandler{
		// private readonly
		db: db,
	}
}

// now to attach the http handler with the construction
// we use golang methods
// in general this fuctions are termed as method function in opps terminology
func (this ListingHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Use the request-scoped context so the database query
	// is cancelled if the client disconnects or the request ends.
	ctx := r.Context()

	rows, err := this.db.QueryContext(
		ctx,
		`SELECT id, title, description, price, city, created_at
		     FROM listings
		     ORDER BY created_at DESC
		     LIMIT 50`)
	if err != nil {
		log.Printf("db.Query: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
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
			log.Printf("rows.Scan: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows.Err: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Serializes the listings slice into JSON
	// and sends it to the client.
	// send error automatically, if caught.
	_ = json.NewEncoder(w).Encode(listings)
}

func (this ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	_, err := this.db.ExecContext(
		ctx,
		`DELETE FROM listings
			WHERE id = $1`, id,
	)
	if err != nil {
		log.Printf("delete: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
