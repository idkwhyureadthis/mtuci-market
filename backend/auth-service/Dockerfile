FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && apk cache clean
ENV CGO_ENABLED=1
RUN go build -o ./auth-service ./cmd/auth-service/main.go


FROM alpine
WORKDIR /app
COPY --from=builder /app/auth-service ./auth-service
EXPOSE 8080
CMD ["./auth-service"]