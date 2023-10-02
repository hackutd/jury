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
func SendJudgeEmail(judge *models.Judge) error {
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
	appName := config.GetEnv("JURY_NAME")
	baseUrl := config.GetEnv("JURY_BASE_URL")

	// Message body
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Jury Judging Platform [%s] \n%s\n\n", appName, mimeHeaders)))

	// Template
	t, _ := template.ParseFiles("email.html")
	t.Execute(&body, struct {
		Name    string
		AppName string
		BaseUrl string
		Code    string
	}{
		Name:    judge.Name,
		AppName: appName,
		BaseUrl: baseUrl,
		Code:    judge.Code,
	})

	// Send email!
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	return err
}
