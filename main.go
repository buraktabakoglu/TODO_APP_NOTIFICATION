package main

import (
	"context"

	"log"

	"github.com/buraktabakoglu/TODO_APP_NOTIFICATION/internal/kafka"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	internal.Consume(context.Background())

}
