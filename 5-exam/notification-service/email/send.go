package email

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"gopkg.in/gomail.v2"
)

func sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "saidazimqilichov@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "saidazimqilichov@gmail.com", "fler dgia ntyz mjoq")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendEmail(to string, name string, email string) (code string, n error) {
	code, err := generateRandomCode()
	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	subject := "----Welcome notification service----"
	body := fmt.Sprintf("Siz joy band qildingiz ismingiz =%s emailingiz=%s", name, email)

	if err := sendEmail(to, subject, body); err != nil {
		log.Printf("Failed to send email: %v", err)
		return "Failed to send", err
	}

	return code, nil

}



func generateRandomCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", n.Int64())
	return code, nil
}
