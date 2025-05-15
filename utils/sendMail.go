package utils

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendExampleEmail () {
	from := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		log.Println("Username and password environment variables must be set")
		return
	}

	to := []string{os.Getenv("TEST_TO_EMAIL")}

	msg := []byte("Subject: Hello from Home Server!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")

	auth := smtp.PlainAuth("Home Server", from, password, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Email sent successfully!")

}

func SendOTP (otp string, toEmail string) error {
	from := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		errMsg := "Username and password environment variables must be set"
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
