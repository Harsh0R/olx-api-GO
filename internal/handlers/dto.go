package handlers

import (
	"fmt"
	"strings"
	"time"
)

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

type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Msg)
}

func (reqData CreateListingObject) Validate() error {
	if strings.TrimSpace(reqData.Title) == "" {
		return &ValidationError{Field: "Title", Msg: "must not empty"}
	}
	return nil
}
