package s3Gateway

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	configApp "github.com/rms-diego/image-processor/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3GatewayService struct{}

type S3GatewayServiceInterface interface {
	Upload(fileHeaders *multipart.FileHeader, file *multipart.File) (*string, *string, error)
}

func NewService() S3GatewayServiceInterface {
	return &S3GatewayService{}
}

func (s *S3GatewayService) Upload(fileHeaders *multipart.FileHeader, file *multipart.File) (*string, *string, error) {
	awsCredentials := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			configApp.Env.AWS_ACCESS_KEY_ID,
			configApp.Env.AWS_SECRET_ACCESS_KEY,
			"",
		),
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(configApp.Env.AWS_REGION),
		config.WithCredentialsProvider(awsCredentials),
	)

	if err != nil {
		return nil, nil, err
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	s3Key := fmt.Sprintf("%v.%v", uuid.New().String(), fileHeaders.Filename)
	s3Res, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(configApp.Env.AWS_S3_BUCKET_NAME),
		Key:    aws.String(s3Key),
		Body:   io.Reader(*file),
	})

	if err != nil {
		return nil, nil, err
	}

	return &s3Res.Location, &s3Key, nil
}
