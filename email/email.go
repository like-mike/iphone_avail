package email

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/like-mike/iphone_avail/config"
	"gopkg.in/gomail.v2"
)

func send(recepient string, subject string, body string, fileName string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Env.Sender)
	m.SetHeader("To", recepient)

	m.SetHeader("Subject", subject)
	m.Attach(fileName)

	if err := config.Env.Smtp.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendRecepients(jobID int64, recepients []string, subject string, body string) error {

	fileName := fmt.Sprintf("%s/%d.txt", config.Env.StaticPath, config.Env.RunID+jobID)
	fmt.Println(fileName)

	err := writeToFile(fileName, body)
	if err != nil {
		return err
	}

	for _, recepient := range recepients {
		log.Printf("%s - Emailing %s...\n", subject, recepient)
		err := send(recepient, subject, body, fileName)
		if err != nil {
			return err
		}
	}

	err = os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(fileName string, data string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}
