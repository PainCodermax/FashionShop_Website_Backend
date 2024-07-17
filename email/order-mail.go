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
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Order Cancellation</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				margin: 0;
				padding: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
			}
			.container {
				background-color: #ffffff;
				padding: 20px 40px;
				border-radius: 8px;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
				text-align: center;
			}
			h1 {
				color: #333;
				font-size: 24px;
			}
			p {
				color: #666;
				font-size: 16px;
			}
			.order-id {
				font-weight: bold;
				color: #000;
			}
			.button {
				display: inline-block;
				margin-top: 20px;
				padding: 10px 20px;
				background-color: #007bff;
				color: #ffffff;
				text-decoration: none;
				border-radius: 5px;
			}
			.button:hover {
				background-color: #0056b3;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>You have successfully cancelled order <span class="order-id">%s</span></h1>
			<p>If you have any questions, please contact our support team.</p>
			<a href="https://fashion-shop-client-379b8.web.app/history/%s" class="button">Go to Orders</a>
		</div>
	</body>
	</html>	
		`
	content = fmt.Sprintf(content, orderID, orderID)

	to := []string{userEmail}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	return err
}
