FROM golang:1.24.4 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG APP

RUN make build-${APP}

FROM debian:stable-slim

WORKDIR /app

COPY --from=build /app/.env .

COPY --from=build /app/build/main .

EXPOSE 8080

#ENV GIN_MODE=release

CMD ["./main"]