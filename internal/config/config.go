package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT   string
	Env    string
	DB_URL string
}

func MustLoad() Config {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is requre!!!")
	}
	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is requre!!!")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("DB_URL is requre!!!")
	}

	return Config{PORT: port, Env: env, DB_URL: dbUrl}
}
