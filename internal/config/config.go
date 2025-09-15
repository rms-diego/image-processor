package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	PORT                  string
	JWT_SECRET            string
	DATABASE_URL          string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	AWS_S3_BUCKET_NAME    string
}

var Env *env

func newEnv() *env {
	return &env{
		PORT:                  os.Getenv("PORT"),
		JWT_SECRET:            os.Getenv("JWT_SECRET"),
		DATABASE_URL:          os.Getenv("DATABASE_URL"),
		AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWS_REGION:            os.Getenv("AWS_REGION"),
		AWS_S3_BUCKET_NAME:    os.Getenv("AWS_S3_BUCKET_NAME"),
	}
}

func Init() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error on load environment variables: %v", err)
	}

	switch {
	case os.Getenv("PORT") == "":
		return fmt.Errorf("PORT variable is not set")

	case os.Getenv("JWT_SECRET") == "":
		return fmt.Errorf("JWT_SECRET variable is not set")

	case os.Getenv("DATABASE_URL") == "":
		return fmt.Errorf("DATABASE_URL variable is not set")

	case os.Getenv("AWS_ACCESS_KEY_ID") == "":
		return fmt.Errorf("AWS_ACCESS_KEY_ID variable is not set")

	case os.Getenv("AWS_SECRET_ACCESS_KEY") == "":
		return fmt.Errorf("AWS_SECRET_ACCESS_KEY variable is not set")

	case os.Getenv("AWS_REGION") == "":
		return fmt.Errorf("AWS_REGION variable is not set")

	case os.Getenv("AWS_S3_BUCKET_NAME") == "":
		return fmt.Errorf("AWS_S3_BUCKET_NAME variable is not set")

	default:
		Env = newEnv()
		return nil
	}
}
