package email

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func SendOPTMail(userEmail string, otp string, isVerify bool) error {
	emailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	emailSenderAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailSenderPassword := os.Getenv("EMAIL_SENDER_PASSWORD")

	sender := NewGmailSender(emailSenderName, emailSenderAddress, emailSenderPassword)
	subject := "YOUR OTP"
	var content string
	if isVerify {
		content = `
			<h1>Welcome to Blawol </h1>
			<p>This is Your OTP %s</p>
			`
		content = fmt.Sprintf(content, otp)
	} else {
		content = `
			<h1>Please enter your OTP to reset password </h1>
			<p>This is Your OTP %s</p>
			`
		content = fmt.Sprintf(content, otp)
	}
	to := []string{userEmail}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	return err
}
