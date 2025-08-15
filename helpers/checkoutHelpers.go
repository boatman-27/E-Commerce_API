package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/go-gomail/gomail"
)

func SendConfirmationEmail(orderId, userEmail string, totalPrice float64, orderDate time.Time) error {
	email := os.Getenv("APP_EMAIL")
	password := os.Getenv("APP_PASSWORD")

	m := gomail.NewMessage()

	m.SetHeader("From", "adhamosman1589@gmail.com") // dev only
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", fmt.Sprintf("Order Confirmation #%s", orderId))

	body := fmt.Sprintf(`Hello,

Thank you for your purchase! Your order #%s has been successfully placed on %s.

Order Details:
- Order ID: %s
- Total Price: $%.2f

We’re processing your order and will notify you once it ships.

If you have any questions, just reply to this email.

Best regards,
Your Friendly Store Team
`, orderId, orderDate.Format("January 2, 2006 at 3:04pm"), orderId, totalPrice)

	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, email, password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
