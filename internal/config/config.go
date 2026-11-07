package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
	Env  string
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

	return Config{PORT: port, Env: env}
}
