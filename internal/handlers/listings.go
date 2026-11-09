package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

type ListingHandler struct {
	db *sql.DB
}

func NewListingHandler(db *sql.DB) *ListingHandler {
	return &ListingHandler{
		db: db,
	}
}

func (lh ListingHandler) List(res http.ResponseWriter, req *http.Request) {
	// req scop context
	ctx := req.Context()
	rows, err := lh.db.QueryContext(ctx, `
		SELECT id, title, description, price, city, created_at 
		FROM listings 
		ORDER BY created_at DESC
		LIMIT 100
		`)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("Request cancelled: %v", err)
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(res, "Request timeout", http.StatusGatewayTimeout)
			return
		}

		log.Printf("listing.query: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
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

func (lh ListingHandler) DeleteList(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := req.PathValue("id")
	fmt.Println("id -> ", id)

	result, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id=$1`, id)
	if err != nil {
		log.Printf("listing.delete: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Printf("listing.rows_affected: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(res, "Listing not found", http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusNoContent)

}
