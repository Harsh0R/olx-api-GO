package main

import (
	"log"
	"net/http"
	"time"

	"github.com/harsh0r/olx-api-GO/internal/config"
)

func main() {

	cfg := config.MustLoad()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /helth", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"Status" : "ok"}`))
	})

	srv := http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	log.Fatalf("Server is running on ... %s" + srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server fail %v", err)
	}
}
