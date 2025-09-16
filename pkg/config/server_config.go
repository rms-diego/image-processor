package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	PORT         string
	JWT_SECRET   string
	DATABASE_URL string
}

var ServerEnv *env

func newEnv() *env {
	return &env{
		PORT:         os.Getenv("PORT"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		DATABASE_URL: os.Getenv("DATABASE_URL"),
	}
}

func InitServerCfg() error {
	godotenv.Load()

	switch {
	case os.Getenv("PORT") == "":
		return fmt.Errorf("PORT variable is not set")

	case os.Getenv("JWT_SECRET") == "":
		return fmt.Errorf("JWT_SECRET variable is not set")

	case os.Getenv("DATABASE_URL") == "":
		return fmt.Errorf("DATABASE_URL variable is not set")

	default:
		ServerEnv = newEnv()
		return nil
	}
}
