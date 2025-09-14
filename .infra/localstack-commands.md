## Export aws keys

```shell
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
```

## CREATE BUCKET IN LOCALSTACK

```
aws --endpoint-url=http://localhost:4566 s3 mb s3://image-processor-bucket
aws --endpoint-url=http://localhost:4566 s3 ls
```

### CREATE SQS QUEUE IN LOCALSTACK

```
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name image-processor-queue
aws --endpoint-url=http://localhost:4566 sqs list-queues
```
