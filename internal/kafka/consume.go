package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	internal "github.com/buraktabakoglu/TODO_APP_NOTIFICATION/internal/controller"
	"github.com/buraktabakoglu/TODO_APP_NOTIFICATION/pkg/models"
	"github.com/segmentio/kafka-go"
)

//birleştir

func Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{os.Getenv("KAFKA_BROKER")},
		Topic:       "email-topic",
		StartOffset: -2,
	})
	topics := []string{"email-topic", "reset_password"}
	for {
		select {
		case <-ctx.Done():
			r.Close()
			return
		default:
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				fmt.Println("could not read message " + err.Error())
				continue
			}
			data := msg.Value
			topic := msg.Topic
			r.CommitMessages(ctx, msg)

			if topic == topics[0] {
				var user models.User
				err = json.Unmarshal(data, &user)
				if err != nil {
					fmt.Println("could not parse userData " + err.Error())
					continue
				}
				token := user.Token
				activationLink := fmt.Sprintf("http://localhost:8080/api/activate/%s", token)
				body := fmt.Sprintf("Your account is now active and your ID is %d. Congrats!", user.ID)
				message := strings.Join([]string{body}, " ")
				go internal.SendEmail(message, user.Email, activationLink)
			} else if topic == topics[1] {
				var resetPassword models.ResetPassword
				err = json.Unmarshal(data, &resetPassword)
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
	}

}
