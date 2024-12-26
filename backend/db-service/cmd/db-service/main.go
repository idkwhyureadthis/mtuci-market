package main

import (
	"auth-service/pkg/handler"
	"os"
)

func main() {
	h := handler.New(os.Getenv("CONN_URL"))

	h.Start(os.Getenv("PORT"))
}
