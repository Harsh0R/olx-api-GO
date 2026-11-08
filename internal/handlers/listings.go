package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price_paisa"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

func List(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		rows, err := db.Query(`
		SELECT id, title, description, price, city, created_at 
		FROM listings 
		ORDER BY created_at DESC
		LIMIT 100
		`)
		if err != nil {
			log.Fatalf("Error Query: %v", err)
			http.Error(res, "internal Error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()
		listings := []listing{}
		for rows.Next() {
			var l listing
			if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
				log.Printf("rows.scan: %v", err)
				http.Error(res, "internal Error", http.StatusInternalServerError)
				return
			}
			listings = append(listings, l)
		}
		if err := rows.Err(); err != nil {
			log.Printf("rows.error: %v", err)
			http.Error(res, "internal error", http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(res).Encode(listings)

	}
}
