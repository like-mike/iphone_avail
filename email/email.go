package email

import (
	"log"

	"github.com/like-mike/iphone_avail/config"
	"gopkg.in/gomail.v2"
)

func send(recepient string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Env.Sender)

	// These have to be sent individually for email-to-text to work
	m.SetHeader("To", recepient)

	//pageURL := "test"
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", `<p>`+body+`</p>`)

	if err := config.Env.Smtp.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendRecepients(recepients []string, subject string, body string) error {
	for _, recepient := range recepients {
		log.Printf("%s - Emailing %s...\n", subject, recepient)
		err := send(recepient, subject, body)
		if err != nil {
			return err
		}
	}
	return nil
}
