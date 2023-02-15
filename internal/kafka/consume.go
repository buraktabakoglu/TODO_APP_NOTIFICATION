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
	"go.uber.org/zap"
)

// consume update

func Consume(ctx context.Context) {
	logger, _ := zap.NewProduction()

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
				logger.Error("could not read message", zap.Error(err))
				continue
			}
			data := msg.Value
			topic := msg.Topic
			r.CommitMessages(ctx, msg)

			if topic == topics[0] {
				var user models.User
				err = json.Unmarshal(data, &user)
				if err != nil {
					logger.Error("could not parse userData", zap.Error(err))
					continue
				}
				token := user.Token
				activationEndpoint := os.Getenv("ACTIVATION_LINK")
				activationLink := activationEndpoint + token
				body := fmt.Sprintf("Your account is now active and your ID is %d. Congrats!", user.ID)
				message := strings.Join([]string{body}, " ")
				go internal.SendEmail(message, user.Email, activationLink)
				err := internal.SendEmail(message, user.Email, activationLink)
				if err != nil {
					logger.Error("could not send email", zap.Error(err))
					continue
				}

			} else if topic == topics[1] {
				var resetPassword models.ResetPassword
				err = json.Unmarshal(data, &resetPassword)
				if err != nil {
					logger.Error("could not parse reset password data", zap.Error(err))
					continue
				}
				token := resetPassword.Token
				email := resetPassword.Email
				resetEndpoint := os.Getenv("RESET_LINK")
				passwordResetLink := resetEndpoint + token
				body := fmt.Sprintf("Please click the following link to reset your password: %s", passwordResetLink)
				message := strings.Join([]string{body}, " ")
				go internal.SendEmail(message, email, "Password reset")
				err := internal.SendEmail(message, resetPassword.Email, passwordResetLink)
				if err != nil {
					logger.Error("could not send email", zap.Error(err))
				}
			}
		}
	}

}
