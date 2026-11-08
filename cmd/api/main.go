package main

import (
	"log"
	"net/http"
	"time"

	"github.com/harsh0r/olx-api-GO/internal/config"
	"github.com/harsh0r/olx-api-GO/internal/db"
	"github.com/harsh0r/olx-api-GO/internal/handlers"
)

func main() {

	cfg := config.MustLoad()

	db, err := db.DbConnect(cfg.DB_URL)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}
	log.Printf("database connected🥳🥳")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /listings", handlers.List(db))
	mux.HandleFunc("DELETE /listings/{id}", handlers.DeleteList(db))

	srv := http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Server is running on ... %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server fail %v", err)
	}
}
