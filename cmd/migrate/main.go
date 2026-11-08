package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/harsh0r/olx-api-GO/internal/config"
)

func main() {

	cfg := config.MustLoad()

	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate <up | down>")
	}
	m, err := migrate.New(
		"file://migrations",
		cfg.DB_URL)

	if err != nil {
		log.Fatalf("migration fail: %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
		log.Println("up called..")
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
		log.Printf("down called..")
	default:
		log.Fatalf("Unknown arg: %s", os.Args[1])
	}

	fmt.Println("migration running....")
}
