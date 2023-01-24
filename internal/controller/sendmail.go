package internal

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(message string, toAddress string, activationLink string) (response bool, err error) {
	fromAddress := os.Getenv("EMAIL")
	fromEmailPassword := os.Getenv("PASSWORD")
	smtpServer := os.Getenv("SMTP_SERVER")
	smptPort := os.Getenv("SMTP_PORT")

	body := fmt.Sprintf("%s\nActivation Link: %s", message, activationLink)
	msg := "From: " + fromAddress + "\n" +
		"To: " + toAddress + "\n" +
		"Subject: Account created!\n\n" +
		body

	var auth = smtp.PlainAuth("", fromAddress, fromEmailPassword, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smptPort, auth, fromAddress, []string{toAddress}, []byte(msg))
	if err == nil {
		return true, nil
	}

	return false, err
}