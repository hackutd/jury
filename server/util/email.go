package util

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"server/config"
	"server/models"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendJudgeEmail sends an email to the judge with their code
func SendJudgeEmail(judge *models.Judge, hostname string) error {
	// If sendgrid API key exists, send email with sendgrid
	sendgridApiKey := config.GetOptEnv("SENDGRID_API_KEY", "")
	if sendgridApiKey != "" {
		return SendgridEmail(sendgridApiKey, judge, hostname)
	}

	// Sender info
	from := config.GetEnv("EMAIL_FROM")
	username := config.GetEnv("EMAIL_USERNAME")
	password := config.GetEnv("EMAIL_PASSWORD")

	// Receiver info
	to := []string{
		judge.Email,
	}

	// SMTP server configuration
	smtpHost := config.GetEnv("EMAIL_HOST")
	smtpPort := config.GetEnv("EMAIL_PORT")

	// Auth
	auth := smtp.PlainAuth("jury", username, password, smtpHost)

	// Other info
	appName := config.GetEnv("REACT_APP_JURY_NAME")

	// Message body
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Jury Judging Platform [%s] \n%s\n\n", appName, mimeHeaders)))

	// Fill Template
	html, err := FillTemplate(judge.Name, hostname, judge.Code, appName)
	if err != nil {
		return err
	}
	body.Write(html)

	// Send email!
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	return err
}

// FillTemplate fills the email template with the given data
func FillTemplate(name string, baseUrl string, code string, appName string) ([]byte, error) {
	var body bytes.Buffer

	t, err := template.ParseFiles("email.html")
	if err != nil {
		return nil, err
	}

	t.Execute(&body, struct {
		Name    string
		BaseUrl string
		Code    string
		AppName string
	}{
		Name:    name,
		BaseUrl: baseUrl,
		Code:    code,
		AppName: appName,
	})

	return body.Bytes(), nil
}

// Send email with Sendgrid
func SendgridEmail(sendgridApiKey string, judge *models.Judge, hostname string) error {
	appName := config.GetEnv("REACT_APP_JURY_NAME")

	from := mail.NewEmail(config.GetEnv("EMAIL_FROM_NAME"), config.GetEnv("EMAIL_FROM"))
	subject := "Jury Judging Platform [" + appName + "]"
	to := mail.NewEmail(judge.Name, judge.Email)

	htmlContent, err := FillTemplate(judge.Name, hostname, judge.Code, appName)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject, to, "", string(htmlContent))
	client := sendgrid.NewSendClient(sendgridApiKey)
	_, err = client.Send(message)
	return err
}
