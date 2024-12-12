package main

import (
	"auth-service/internal/endpoint"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	secret := []byte(os.Getenv("SECRET_KEY"))
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port not provided")
	}
	srv := endpoint.New(secret)
	srv.Start(port)
}
