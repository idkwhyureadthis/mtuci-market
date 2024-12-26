FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && apk cache clean
ENV CGO_ENABLED=1
RUN go build -o ./api-service ./cmd/api-service/main.go


FROM alpine
WORKDIR /app
COPY --from=builder /app/api-service ./api-service
EXPOSE 8082
CMD ["./api-service"]