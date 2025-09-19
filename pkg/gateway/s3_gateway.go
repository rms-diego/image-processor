package gateway

import (
	"bytes"
	"context"

	configApp "github.com/rms-diego/image-processor/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Gateway struct {
	client *s3.Client
	ctx    context.Context
}

var S3Gateway *s3Gateway

type S3GatewayInterface interface {
	Upload(s3Key *string, fileBytes *[]byte) (*string, error)
	RemoveObject(s3Key *string) error
	GetObject(s3Key *string) (*s3.GetObjectOutput, error)
}

func newService() *s3Gateway {
	client := s3.NewFromConfig(configApp.GatewayCfg.AWS_CFG)

	return &s3Gateway{client: client, ctx: context.TODO()}
}

func InitS3() {
	S3Gateway = newService()
}

func (s *s3Gateway) Upload(s3Key *string, fileBytes *[]byte) (*string, error) {
	uploader := manager.NewUploader(s.client)

	s3Res, err := uploader.Upload(s.ctx, &s3.PutObjectInput{
		Bucket: aws.String(configApp.GatewayCfg.AWS_S3_BUCKET_NAME),
		Key:    aws.String(*s3Key),
		Body:   bytes.NewReader(*fileBytes),
	})

	if err != nil {
		return nil, err
	}

	return &s3Res.Location, nil
}

func (s *s3Gateway) RemoveObject(s3Key *string) error {
	_, err := s.client.DeleteObject(s.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(configApp.GatewayCfg.AWS_S3_BUCKET_NAME),
		Key:    s3Key,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *s3Gateway) GetObject(s3Key *string) (*s3.GetObjectOutput, error) {
	output, err := s.client.GetObject(s.ctx, &s3.GetObjectInput{
		Bucket: aws.String(configApp.GatewayCfg.AWS_S3_BUCKET_NAME),
		Key:    s3Key,
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
