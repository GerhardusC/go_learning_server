package utils

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

func SendExampleEmail () {
	from := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		log.Println("Username and password environment variables must be set")
		return
	}

	to := []string{os.Getenv("TEST_TO_EMAIL")}

	now := time.Now()

	msg := []byte("Subject: Home Server has started!\r\n" +
		"\r\n" +
		"Started at:\r\n" + now.String())

	auth := smtp.PlainAuth("Home Server", from, password, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Email sent successfully!")
}

func SendOTPEmail (otp string, toEmail string) error {
	from := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		errMsg := "EMAIL_USERNAME and EMAIL_password environment variables must be set"
		log.Println(errMsg)
		return errors.New(errMsg)
	}

	to := []string{toEmail}

	msg := fmt.Appendf(nil, `Subject: Login OTP

		OTP: %s
	`, otp)

	auth := smtp.PlainAuth("Home Server", from, password, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		log.Println("Something went wrong while sending email:", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}
