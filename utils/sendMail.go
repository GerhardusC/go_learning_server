package utils

import (
	"os"
	"net/smtp"
	"log"
	"fmt"
)

func SendExampleEmail () {
	from := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		log.Println("Username and password environment variables must be set")
		return
	}

	to := []string{"test.test@gmail.com"}

	msg := []byte("Subject: Hello from Home Server!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")

	auth := smtp.PlainAuth("Home Server", from, password, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")

}
