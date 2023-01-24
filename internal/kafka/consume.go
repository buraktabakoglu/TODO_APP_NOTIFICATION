package Kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	internal "github.com/buraktabakoglu/TODO_APP_NOTIFICATION.git/internal/controller"
	"github.com/buraktabakoglu/TODO_APP_NOTIFICATION.git/pkg/models"
	"github.com/segmentio/kafka-go"
)

func Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{os.Getenv("KAFKA_BROKER")},
		Topic:       "email-topic",
		StartOffset: kafka.LastOffset,
	})

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		userData := msg.Value

		r.CommitMessages(ctx, msg)

		var user models.User

		fmt.Println("gönderiliyor...")

		err = json.Unmarshal(userData, &user)
		if err != nil {
			panic("could not parse userData " + err.Error())
		}
		token := user.Token

		fmt.Println(token)
		activationLink := fmt.Sprintf("http://localhost:8080/api/activate/%s", token)
		body := fmt.Sprintf("Your account is now active and your ID is %d. Congrats!", user.ID)
		message := strings.Join([]string{body}, " ")
		go internal.SendEmail(message, user.Email, activationLink)

	}
}

func ConsumePassword(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{os.Getenv("KAFKA_BROKER")},
		Topic:       "reset_password",
		StartOffset: kafka.LastOffset,
	})

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Println("could not read message " + err.Error())
			continue
		}
		resetPasswordData := msg.Value

		r.CommitMessages(ctx, msg)

		var resetPassword models.ResetPassword

		fmt.Println("Sending password reset link...")

		err = json.Unmarshal(resetPasswordData, &resetPassword)
		if err != nil {
			fmt.Println("could not parse reset password data " + err.Error())
			continue
		}
		token := resetPassword.Token
		email := resetPassword.Email

		passwordResetLink := fmt.Sprintf("http://localhost:8080/api/password/reset/%s", token)
		body := fmt.Sprintf("Please click the following link to reset your password: %s", passwordResetLink)
		message := strings.Join([]string{body}, " ")
		go internal.SendEmail(message, email, "Password reset")

	}
}
