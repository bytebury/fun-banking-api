package mailers

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/mail.v2"
)

/**
 * Sends an email.
 */
func sendEmail(to, subject, templateName string, data any) error {
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))

	email := mail.NewMessage()
	email.SetHeader("From", username)
	email.SetHeader("To", to)
	email.SetHeader("Subject", subject)

	if err := createPlainTextBody(email, data, fmt.Sprintf("%s.txt", templateName)); err != nil {
		return err
	}

	if err := createHtmlBody(email, data, fmt.Sprintf("%s.html", templateName)); err != nil {
		return err
	}

	dialer := mail.NewDialer(host, port, username, password)

	if err := dialer.DialAndSend(email); err != nil {
		fmt.Println("📪 Unable to send email to: ", to)
		return err
	}

	fmt.Println("📬 Sent an email to: ", to)

	return nil
}

func createHtmlBody(email *mail.Message, data any, templateName string) error {
	return createEmailBody(email, data, templateName, "text/html")
}

func createPlainTextBody(email *mail.Message, data any, templateName string) error {
	return createEmailBody(email, data, templateName, "text/plain")
}

func createEmailBody(email *mail.Message, data any, templateName, templateType string) error {
	templ, err := template.ParseFiles(fmt.Sprintf("mailers/templates/%s", templateName))

	if err != nil {
		return err
	}

	var emailContent string

	buffer := &bytes.Buffer{}
	if err := templ.Execute(buffer, &data); err != nil {
		return err
	}

	emailContent = buffer.String()
	email.SetBody(templateType, emailContent)

	return nil
}
