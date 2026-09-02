package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

func SendApiKeyEmail(toEmail, userName, apiKey string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.office365.com"
	}
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort := 587
	if p, err := strconv.Atoi(smtpPortStr); err == nil && p > 0 {
		smtpPort = p
	}

	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		smtpUser = "ai@safexnet.com"
	}
	smtpPass := os.Getenv("SMTP_PASSWORD")
	if smtpPass == "" {
		smtpPass = "River!Cloud@82Forest"
	}
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "ai@safexnet.com"
	}
	fromName := os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		fromName = "SafexNext"
	}

	subject := "Welcome to SafexNext - Your API Key"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; background-color: #f4f6f8; padding: 20px;">
  <div style="max-width: 600px; margin: 0 auto; background: #ffffff; padding: 30px; border-radius: 8px;">
    <h2 style="color: #1a202c;">Welcome to SafexNext, %s!</h2>
    <p>Your account has been created successfully.</p>
    <p>Below is your personal <strong>API Key</strong> required for authentication:</p>
    <div style="background: #edf2f7; padding: 15px; border-radius: 6px; font-family: monospace; font-size: 18px; font-weight: bold; text-align: center; margin: 20px 0;">
      %s
    </div>
    <p style="color: #718096; font-size: 13px;">Keep this key secure and do not share it with anyone.</p>
  </div>
</body>
</html>`, userName, apiKey)

	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", fromName, fromEmail)
	header["To"] = toEmail
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Establish TLS connection for Office 365 (Port 587 STARTTLS)
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpHost, smtpPort))
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer c.Close()

	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = c.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err = c.Mail(fromEmail); err != nil {
		return fmt.Errorf("SMTP mail from failed: %w", err)
	}

	if err = c.Rcpt(toEmail); err != nil {
		return fmt.Errorf("SMTP rcpt to failed: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("SMTP data failed: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message body: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return c.Quit()
}
