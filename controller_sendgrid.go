package main

import (
	"fmt"
	"log"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendgrid_email_credentials(username string, useremail string, userpassword string) {
	from := mail.NewEmail("Cavalry Tactics Inc.", "noreply@cavalrytactics.com")
	subject := "Welcome to Cavalry Tactics!"
	to := mail.NewEmail(username, useremail)
	plainTextContent := "The password to your account is: " + userpassword
	htmlContent := "<strong>The password to your account is: " + userpassword + "</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendgridkey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
