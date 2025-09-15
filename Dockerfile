FROM golang:1.24-alpine3.21 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o build/main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/.env .

COPY --from=build /app/build/main .

EXPOSE 8080

ENV GIN_MODE=release

CMD ["./main"]