package email

import (
	"fmt"
	"os"
)

func ConfirmOrder(userEmail string, orderID string) error {
	emailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	emailSenderAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailSenderPassword := os.Getenv("EMAIL_SENDER_PASSWORD")

	sender := NewGmailSender(emailSenderName, emailSenderAddress, emailSenderPassword)
	subject := "CONFIRM ORDER ON BLAWOL"
	var content string

	content = `
		<h1>You have successfully placed order %s</h1>
		`
	content = fmt.Sprintf(content, orderID)

	to := []string{userEmail}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	return err
}

func CancelOrder(userEmail string, orderID string) error {
	emailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	emailSenderAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailSenderPassword := os.Getenv("EMAIL_SENDER_PASSWORD")

	sender := NewGmailSender(emailSenderName, emailSenderAddress, emailSenderPassword)
	subject := "CANCEL ORDER ON BLAWOL"
	var content string

	content = `
		<h1>You have successfully cancelled order %s</h1>
		`
	content = fmt.Sprintf(content, orderID)

	to := []string{userEmail}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	return err
}

