package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type awsEnv struct {
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	AWS_S3_BUCKET_NAME    string
	AWS_SQS_URL           string
}

var AwsEnv *awsEnv

func newAwsEnv() *awsEnv {
	return &awsEnv{
		AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWS_REGION:            os.Getenv("AWS_REGION"),
		AWS_S3_BUCKET_NAME:    os.Getenv("AWS_S3_BUCKET_NAME"),
		AWS_SQS_URL:           os.Getenv("AWS_SQS_URL"),
	}
}

func InitAwsCfg() error {
	godotenv.Load()

	switch {
	case os.Getenv("AWS_ACCESS_KEY_ID") == "":
		return fmt.Errorf("AWS_ACCESS_KEY_ID variable is not set")

	case os.Getenv("AWS_SECRET_ACCESS_KEY") == "":
		return fmt.Errorf("AWS_SECRET_ACCESS_KEY variable is not set")

	case os.Getenv("AWS_REGION") == "":
		return fmt.Errorf("AWS_REGION variable is not set")

	case os.Getenv("AWS_S3_BUCKET_NAME") == "":
		return fmt.Errorf("AWS_S3_BUCKET_NAME variable is not set")

	case os.Getenv("AWS_SQS_URL") == "":
		return fmt.Errorf("AWS_SQS_URL variable is not set")

	default:
		AwsEnv = newAwsEnv()
		return nil
	}
}
