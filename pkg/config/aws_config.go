package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/joho/godotenv"
)

type awsCfg struct {
	AWS_S3_BUCKET_NAME string
	AWS_SQS_URL        string
	AWS_CFG            aws.Config
}

var AwsCfg *awsCfg

func newAwsCfg() (*awsCfg, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	awsCredentials := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		),
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(awsCredentials),
	)

	if err != nil {
		return nil, err
	}

	return &awsCfg{
		AWS_SQS_URL:        os.Getenv("AWS_SQS_URL"),
		AWS_S3_BUCKET_NAME: os.Getenv("AWS_S3_BUCKET_NAME"),
		AWS_CFG:            cfg,
	}, nil
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
		cfg, err := newAwsCfg()

		if err != nil {
			return err
		}

		AwsCfg = cfg
		return nil
	}
}
