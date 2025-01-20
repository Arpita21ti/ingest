package utils

import (
	"fmt"
	"net/smtp"
)

// SendResetEmail sends a password reset email to the user
func SendResetEmail(to, tokenString string) error {
	// SMTP server configuration
	const (
		smtpServer = "smtp.example.com"    // Replace with your SMTP server
		port       = "587"                 // Replace with your SMTP port
		username   = "your-smtp-username"  // Replace with your SMTP username
		password   = "your-smtp-password"  // Replace with your SMTP password
		from       = "noreply@example.com" // Replace with your official "From" email
	)

	// Update with the reset url.
	resetLink := fmt.Sprintf("https://example.com/reset-password?token=%s", tokenString)

	// Subject and body of the email
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		Hello,

		We received a request to reset your password. Please click the link below to reset your password:
		
		%s

		If you did not request this, you can ignore this email.

		Thank you,
		Your Team
	`, resetLink)

	// Prepare the email message
	message := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, subject, body,
	))

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpServer, port),
		smtp.PlainAuth("", username, password, smtpServer),
		from, []string{to}, message,
	)
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %v", to, err)
	}

	return nil
}
