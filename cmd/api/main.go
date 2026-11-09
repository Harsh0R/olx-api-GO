package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/harsh0r/olx-api-GO/internal/config"
	"github.com/harsh0r/olx-api-GO/internal/db"
	"github.com/harsh0r/olx-api-GO/internal/handlers"
	"github.com/harsh0r/olx-api-GO/internal/middleware"
)

func main() {

	cfg := config.MustLoad()

	db, err := db.DbConnect(cfg.DB_URL)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	log.Printf("database connected😘😘")
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	lh := handlers.NewListingHandler(db, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.DeleteList)

	handler := middleware.RequestId(mux)

	srv := http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Server is running on ... %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server fail %v", err)
	}
}
