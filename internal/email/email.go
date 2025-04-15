package emailing

import (
	"fmt"

	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

func SendOTPEmail(email, otp, subject string) (string, error) {
	expirationTime := "10 minutes"
	body := fmt.Sprintf(`
	<html>
	<body>
	<img src="https://res.cloudinary.com/dhlfvmmd6/image/upload/v1741631529/event_images/img_1741631526728306000.png" 
	alt="BildWerk Logo" style="margin-bottom: 20px;">
	<p>Your One-Time Password (OTP) is: <strong>%s</strong>.</p>
	<p>This code is valid for %s.</p>
	<p>If you did not request this, please ignore this email.</p>
	<p>Thank you for using BildWerk!</p>
	</body>
	</html>`, otp, expirationTime)

	if err := SendEmailWithRetry(email, subject, body, true, 3); err != nil {
		log.Error().Err(err).Msg("Failed to send OTP email")
		return "", err
	}

	log.Info().Str("email", email).Msg("OTP email sent successfully")
	return "OTP sent successfully", nil
}

func SendWelcomeEmail(email string) (string, error) {
	subject := "Welcome to BildWerk!"
	body := `
	<html>
	<body>
	<img src="https://res.cloudinary.com/dhlfvmmd6/image/upload/v1741631529/event_images/img_1741631526728306000.png" 
	alt="BildWerk Logo" style="margin-bottom: 20px;">
	<p>Hi,</p>
	<p>Thank you for registering with us! We are excited to have you on board.</p>
	<p>BildWerk is dedicated to providing you with the best experience. Here are a few things you can do:</p>
	<ul>
		<li>Explore our features and services.</li>
		<li>Set up your profile to get started.</li>
		<li>Contact our support team if you have any questions.</li>
	</ul>
	<p>If you need assistance, feel free to reach out to our support team at any time.</p>
	<p>Thank you for joining BildWerk! We look forward to serving you.</p>
	<p>Best regards,<br>The BildWerk Team</p>
	</body>
	</html>`

	if err := SendEmailWithRetry(email, subject, body, true, 3); err != nil {
		log.Error().Err(err).Msg("Failed to send welcome email")
		return "", err
	}

	log.Info().Str("email", email).Msg("Welcome email sent successfully")
	return "Welcome email sent successfully", nil
}

func SendEmailWithRetry(to, subject, body string, isHTML bool, maxRetries int) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	intValuePort, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	d := gomail.NewDialer(host, intValuePort, username, password)

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = sendEmail(d, from, to, subject, body, isHTML)
		if err == nil {
			return nil
		}
		log.Warn().Str("addr", host).Int("attempt", attempt+1).Err(err).Msg("Failed to send email, retrying...")
		time.Sleep(time.Duration(attempt+1) * time.Second) // Exponential backoff
	}
	return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, err)
}

func sendEmail(d *gomail.Dialer, from, to, subject, body string, isHTML bool) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	if isHTML {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
