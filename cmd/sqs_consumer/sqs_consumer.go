package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rms-diego/image-processor/internal/database"
	imagerepository "github.com/rms-diego/image-processor/internal/modules/image/image_repository"
	imageservice "github.com/rms-diego/image-processor/internal/modules/image/image_service"
	"github.com/rms-diego/image-processor/internal/validations"
	"github.com/rms-diego/image-processor/pkg/config"
	"github.com/rms-diego/image-processor/pkg/gateway"
)

func main() {

	if err := config.InitGatewayCfg(); err != nil {
		panic(err)
	}

	if err := database.Init(config.GatewayCfg.DATABASE_URL); err != nil {
		panic(err)
	}

	gateway.InitSQS()
	gateway.InitS3()

	ir := imagerepository.NewImageRepository(database.DB)
	is := imageservice.NewService(gateway.S3Gateway, gateway.SqsGateway, ir)

	for {
		sqsMessages, err := gateway.SqsGateway.GetMessages()
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			continue
		}

		if err != nil {
			fmt.Println("Error getting messages from SQS:", err)
			continue
		}

		for _, msg := range sqsMessages {
			fmt.Println("Processing message ID:", *msg.MessageId)

			if err := processMessages(msg, is); err != nil {
				break // TODO: Implement DLQ (Dead Letter Queue)
			}

			if err := gateway.SqsGateway.RemoveMessage(msg.ReceiptHandle); err != nil {
				fmt.Println("Error removing message from SQS:", err)
				break
			}

			continue
		}
	}
}

func processMessages(msg types.Message, is imageservice.ImageServiceInterface) error {
	var formatMessage validations.TransformMessageQueue
	if err := json.Unmarshal([]byte(*msg.Body), &formatMessage); err != nil {
		fmt.Println("Error unmarshalling message body:", err)
		return err
	}

	s3Object, err := gateway.S3Gateway.GetObject(&formatMessage.S3Key)
	if err != nil {
		fmt.Println("Error getting object from S3:", err)
		return err
	}

	file := s3Object.Body
	defer file.Close()

	fileBuffer, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	fmt.Println("Successfully retrieved object from S3")
	if err := is.ProcessImage(&fileBuffer, &formatMessage); err != nil {
		fmt.Println("Error processing image:", err)
		return err
	}

	return nil
}
