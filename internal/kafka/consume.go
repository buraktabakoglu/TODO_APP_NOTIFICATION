package kafka

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

const (
	ActivationKey = "register"
	ResetKey      = "reset_password"
)

var logger *zap.Logger

func init() {

	var err error

	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
}

func Consume(ctx context.Context) {

	topic := os.Getenv("TOPIC")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{os.Getenv("KAFKA_BROKER")},
		Topic:       topic,
		StartOffset: kafka.LastOffset,
	})
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
			key := string(msg.Key)

			if key == ActivationKey {
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
				go func() {
					if err := internal.SendEmail(message, user.Email, activationLink); err != nil {
						logger.Error("could not send email", zap.Error(err))
					}
					r.CommitMessages(ctx, msg)
				}()

			} else if key == ResetKey {
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
				go func() {
					if err := internal.SendEmail(message, email, "Password reset"); err != nil {
						logger.Error("could not send email", zap.Error(err))
					}
					r.CommitMessages(ctx, msg)
				}()
			}
		}
	}

}
