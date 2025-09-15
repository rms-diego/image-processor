# image-processor

build image

```shell
docker build -t image-processor-golang .
```

run container

```shell
docker run --env-file .env -p 8080:8080 image-processor-golang
```
