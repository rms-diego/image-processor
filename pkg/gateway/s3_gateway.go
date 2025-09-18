package gateway

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	configApp "github.com/rms-diego/image-processor/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type s3Gateway struct {
	client *s3.Client
}

var S3Gateway *s3Gateway

type S3GatewayInterface interface {
	Upload(fileHeaders *multipart.FileHeader, file *multipart.File) (*string, *string, error)
	GetObject(s3Key *string) (*s3.GetObjectOutput, error)
}

func newService() *s3Gateway {
	client := s3.NewFromConfig(configApp.GatewayCfg.AWS_CFG)

	return &s3Gateway{client: client}
}

func InitS3() {
	S3Gateway = newService()
}

func (s *s3Gateway) Upload(fileHeaders *multipart.FileHeader, file *multipart.File) (*string, *string, error) {
	uploader := manager.NewUploader(s.client)

	s3Key := fmt.Sprintf("%v.%v", uuid.New().String(), fileHeaders.Filename)
	s3Res, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(configApp.GatewayCfg.AWS_S3_BUCKET_NAME),
		Key:    aws.String(s3Key),
		Body:   io.Reader(*file),
	})

	if err != nil {
		return nil, nil, err
	}

	return &s3Res.Location, &s3Key, nil
}

func (s *s3Gateway) GetObject(s3Key *string) (*s3.GetObjectOutput, error) {
	output, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(configApp.GatewayCfg.AWS_S3_BUCKET_NAME),
		Key:    s3Key,
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
