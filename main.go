package main

import (
	"context"

	"log"

	"github.com/buraktabakoglu/TODO_APP_NOTIFICATION.git/internal/Kafka"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	Kafka.Consume(context.Background())
	Kafka.ConsumePassword(context.Background())

}
