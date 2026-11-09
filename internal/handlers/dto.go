package handlers

import "time"

// DTO : data transfer object
type CreateListingObject struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price_paisa"`
	City        string `json:"city"`
}
type CreateListingResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
