package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/harsh0r/olx-api-GO/internal/httpx"
	"github.com/harsh0r/olx-api-GO/internal/middleware"
)

type Listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price_paisa"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

// GET /listings
func (lh *ListingHandler) List(
	res http.ResponseWriter,
	req *http.Request,
) {
	ctx := req.Context()

	rows, err := lh.db.QueryContext(ctx, `
		SELECT id, title, description, price, city, created_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT 100
	`)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			lh.logger.Warn("Listing request cancelled",
				"method", req.Method,
				"path", req.URL.Path,
				"error", err,
			)
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			lh.logger.Warn("Listing query timed out",
				"method", req.Method,
				"path", req.URL.Path,
				"error", err,
			)

			http.Error(res, "Request timeout", http.StatusGatewayTimeout)
			return
		}

		lh.logger.Error("Failed to fetch listings",
			"method", req.Method,
			"path", req.URL.Path,
			"error", err,
		)

		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	listings := make([]Listing, 0)

	for rows.Next() {
		var l Listing

		if err := rows.Scan(
			&l.ID,
			&l.Title,
			&l.Description,
			&l.Price,
			&l.City,
			&l.CreatedAt,
		); err != nil {
			lh.logger.Error("Failed to scan listing row",
				"error", err,
			)

			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		lh.logger.Error("Error while iterating listing rows",
			"error", err,
		)

		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(res).Encode(listings); err != nil {
		lh.logger.Error("Failed to encode listings response",
			"error", err,
		)
		return
	}
}

// DELETE /listings/{id}
func (lh *ListingHandler) DeleteList(
	res http.ResponseWriter,
	req *http.Request,
) {
	ctx := req.Context()

	requestId := middleware.RequestIdFromContext(ctx)
	// requestId := ctx.Value("requestCtxId").(string)
	id := req.PathValue("id")

	if id == "" {
		lh.logger.Warn("Delete listing request with missing ID",
			"method", req.Method,
			"path", req.URL.Path,
		)

		http.Error(res, "Listing ID is required", http.StatusBadRequest)
		return
	}

	result, err := lh.db.ExecContext(
		ctx,
		`DELETE FROM listings WHERE id = $1`,
		id,
	)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			lh.logger.Warn("Delete listing request cancelled",
				"listing_id", id,
				"error", err,
				"request_id", requestId,
			)
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			lh.logger.Warn("Delete listing query timed out",
				"listing_id", id,
				"error", err,
				"request_id", requestId,
			)

			http.Error(res, "Request timeout", http.StatusGatewayTimeout)
			return
		}

		lh.logger.Error("Failed to delete listing",
			"listing_id", id,
			"error", err,
			"request_id", requestId,
		)

		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()

	if err != nil {
		lh.logger.Error("Failed to get affected rows",
			"listing_id", id,
			"error", err,
			"request_id", requestId,
		)

		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		lh.logger.Warn("Listing not found",
			"listing_id", id,
			"request_id", requestId,
		)
		httpx.Error(res, http.StatusInternalServerError, "somathing went wrong", "internal_error")
		// http.Error(res, "Listing not found", http.StatusNotFound)
		return
	}

	lh.logger.Info("Listing deleted successfully",
		"listing_id", id,
	)

	res.WriteHeader(http.StatusNoContent)
}
