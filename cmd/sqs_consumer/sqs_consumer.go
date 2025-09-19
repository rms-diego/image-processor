package main

import (
	"encoding/json"
	"fmt"
	"time"

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
		if err != nil {
			fmt.Println("Error getting messages from SQS:", err)
			time.Sleep(800 * time.Millisecond)
			continue
		}

		if len(sqsMessages) == 0 {
			time.Sleep(800 * time.Millisecond)
			continue
		}

		for _, msg := range sqsMessages {
			fmt.Println("Processing message ID:", *msg.MessageId)
			fmt.Println("Message Body:", *msg.Body)

			if err := processMessages(msg, is); err != nil {
				return // TODO: IMPLEMENTS DLQ
			}
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

	fmt.Println("Successfully retrieved object from S3:", s3Object)
	if err := is.ProcessImage(&s3Object.Body, &formatMessage.Payload); err != nil {
		fmt.Println("Error processing image:", err)
		return err
	}

	return nil
}
