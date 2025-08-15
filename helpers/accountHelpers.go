package helpers

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-gomail/gomail"
	"github.com/jmoiron/sqlx"
)

func IsValidEmail(email string) bool {
	// Simple regex pattern for email validation
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func IsEmailAvailable(db *sqlx.DB, email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil // true = available
}

func IsUserIdAvailable(db *sqlx.DB, userId string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil // true = available
}

func SendVerificationEmail(token, userEmail string) error {
	email := os.Getenv("APP_EMAIL")
	password := os.Getenv("APP_PASSWORD")

	m := gomail.NewMessage()

	m.SetHeader("From", "adhamosman1589@gmail.com") // dev only
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Email Verification")

	verificationLink := fmt.Sprintf("http://localhost:3000/verify?verificationToken=%s", token)

	body := fmt.Sprintf("Please verify your account by clicking the following link:\n\n%s", verificationLink)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, email, password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
