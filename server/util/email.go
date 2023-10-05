package util

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"server/config"
	"server/models"
)

// SendJudgeEmail sends an email to the judge with their code
func SendJudgeEmail(judge *models.Judge, hostname string) error {
	// Sender info
	from := config.GetEnv("EMAIL_FROM")
	password := config.GetEnv("EMAIL_PASSWORD")

	// Receiver info
	to := []string{
		judge.Email,
	}

	// SMTP server configuration
	smtpHost := config.GetEnv("EMAIL_HOST")
	smtpPort := config.GetEnv("EMAIL_PORT")

	// Auth
	auth := smtp.PlainAuth("jury", from, password, smtpHost)

	// Other info
	appName := config.GetEnv("REACT_APP_JURY_NAME")

	// Message body
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Jury Judging Platform [%s] \n%s\n\n", appName, mimeHeaders)))

	// Template
	t, _ := template.ParseFiles("email.html")
	t.Execute(&body, struct {
		Name    string
		BaseUrl string
		Code    string
		AppName string
	}{
		Name:    judge.Name,
		BaseUrl: hostname,
		Code:    judge.Code,
		AppName: appName,
	})

	// Send email!
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	return err
}
