# image-processor

## Tech Stack

[![My Skills](https://skillicons.dev/icons?i=go,postgres,docker,aws)](https://skillicons.dev)

## About

Asynchronous Image Processing with AWS SQS and S3. This project showcases a scalable architecture for asynchronous image processing using AWS services.

🚀 Features

- Upload images through the application.
- Queue processing tasks with Amazon SQS, ensuring decoupling between producers and consumers.
- Asynchronous workers handle image operations (e.g., resizing, compression...).
- Retry dead letter in DLQ queue
- Store processed images in Amazon S3.

Verify api documentation in `/docs`

## Running local

1. **Install Dependencies**: `go mod tidy`
2. **Environment variables**: Copy `.env.example` to a new `.env`
3. **Create postgres entities**: run queries in `.infra/init.sql`
4. **Run Api**: `make run-api`
5. **Run SQS consumers**: `make run-sqs-consumer`

## Running with docker

1. **Environment variables**: Copy `.env.example` to a new `.env`
2. **Create containers**: `docker-compose up -d`

## Contribute

1. **Clone project**: `git clone https://github.com/rms-diego/image-processor.git`
2. **Create feature/branch**: `git checkout -b feature/NAME`

## License

This software is available under the following licenses:

- [MIT](https://rem.mit-license.org)
- [Challenger link](https://roadmap.sh/projects/image-processing-service)

<!-- # image-processor

build image

```shell
docker build --build-arg APP=api -t image-processor-golang .
```

run container

```shell
docker run --env-file .env -p 8080:8080 image-processor-golang
``` -->
