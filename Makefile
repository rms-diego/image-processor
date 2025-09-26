run-api: 
	go run cmd/main.go

run-api-dev:
	air

build-api:
	go build -o ./build/main ./cmd/api/main.go

run-sqs-consumer:
	go run cmd/sqs_consumer/sqs_consumer.go

build-sqs-consumer:
	go build -o ./build/main ./cmd/sqs_consumer/sqs_consumer.go