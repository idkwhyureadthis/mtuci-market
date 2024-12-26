FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && apk cache clean
ENV CGO_ENABLED=1
RUN go build -o ./db-service ./cmd/db-service/main.go


FROM alpine
WORKDIR /app
COPY --from=builder /app/db-service ./db-service
COPY --from=builder /app/internal/migrations ./internal/migrations
EXPOSE 8081
CMD ["./db-service"]