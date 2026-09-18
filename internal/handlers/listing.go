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

// closure factory
func List(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(
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

		listings := []listing{}
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
}
