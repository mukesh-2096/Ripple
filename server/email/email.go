package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends a transactional email. Falls back to logging if SMTP is unconfigured.
func SendEmail(to string, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if from == "" {
		from = "no-reply@ripple.dev"
	}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)

	// If SMTP parameters are missing, fallback to logging
	if smtpHost == "" || smtpPort == "" {
		log.Printf("[EMAIL SIMULATION] Sending email to: %s\nSubject: %s\nBody: %s", to, subject, body)
		// Write to a local file for proof of email generation without SMTP
		f, err := os.OpenFile("sent_emails.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString(fmt.Sprintf("=== EMAIL SENT AT %s ===\nTo: %s\nSubject: %s\nBody:\n%s\n=========================\n\n", 
				fmt.Sprint(os.Getenv("LOCAL_TIME")), to, subject, body))
		}
		return nil
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("SMTP Error: Failed to send email to %s: %v", to, err)
		return err
	}

	return nil
}
