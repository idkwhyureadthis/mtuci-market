package main

import (
	"api-service/internal/endpoint"
	"os"
)

func main() {
	port := ":" + os.Getenv("PORT")

	e := endpoint.New()

	e.Start(port)
}
