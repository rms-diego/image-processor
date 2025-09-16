package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	configApp "github.com/rms-diego/image-processor/pkg/config"
)

type sqsGateway struct {
	sqsClient *sqs.Client
	sqsUrl    string
}

type SqsGatewayServiceInterface interface {
	SendMessage(messageBody string) error
	GetMessages() ([]types.Message, error)
}

var SqsGateway *sqsGateway

func InitSQS() error {
	ctx := context.TODO()
	awsCredentials := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			configApp.AwsEnv.AWS_ACCESS_KEY_ID,
			configApp.AwsEnv.AWS_SECRET_ACCESS_KEY,
			"",
		),
	)

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(configApp.AwsEnv.AWS_REGION),
		config.WithCredentialsProvider(awsCredentials),
	)

	if err != nil {
		return err
	}

	client := sqs.NewFromConfig(cfg)
	SqsGateway = newSqsGateway(client)
	return nil
}

func newSqsGateway(client *sqs.Client) *sqsGateway {
	return &sqsGateway{sqsClient: client, sqsUrl: configApp.AwsEnv.AWS_SQS_URL}
}

func (g *sqsGateway) SendMessage(messageBody string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := g.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &g.sqsUrl,
		MessageBody: &messageBody,
	})

	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	return nil
}

func (g *sqsGateway) GetMessages() ([]types.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output, err := g.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &g.sqsUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     5,
		VisibilityTimeout:   30,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to receive messages from SQS: %w", err)
	}

	return output.Messages, nil
}
