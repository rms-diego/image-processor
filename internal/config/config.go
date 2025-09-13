package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	PORT        string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

var Env *env

func newEnv() *env {
	return &env{
		PORT:        os.Getenv("PORT"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
	}
}

func Init() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error on load environment variables: %v", err)
	}

	switch {
	case os.Getenv("PORT") == "":
		return fmt.Errorf("PORT variable is not set")

	case os.Getenv("DB_HOST") == "":
		return fmt.Errorf("DB_HOST variable is not set")

	case os.Getenv("DB_PORT") == "":
		return fmt.Errorf("DB_PORT variable is not set")

	case os.Getenv("DB_USER") == "":
		return fmt.Errorf("DB_USER variable is not set")

	case os.Getenv("DB_PASSWORD") == "":
		return fmt.Errorf("DB_PASSWORD variable is not set")

	case os.Getenv("DB_NAME") == "":
		return fmt.Errorf("DB_NAME variable is not set")

	default:
		Env = newEnv()
		return nil
	}
}
