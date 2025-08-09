FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /migrator ./cmd/migrator
RUN CGO_ENABLED=0 GOOS=linux go build -o /subscriptions ./cmd/subscriptions

FROM alpine:3.20

WORKDIR /app

ENV CONFIG_PATH=./configs/env.yaml

COPY --from=builder /migrator /migrator
COPY --from=builder /subscriptions /subscriptions
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

ENTRYPOINT []
CMD ["/subscriptions"]
