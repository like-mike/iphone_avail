package config

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"strings"
	"time"

	"fmt"
	"gopkg.in/gomail.v2"
)

var (
	Env Config
)

type Config struct {
	Ctx             context.Context
	Sender          string
	Recepients      []string
	ErrorRecepients []string
	SmtpHost        string
	SmtpPass        string
	SmtpHeader      string
	Smtp            *gomail.Dialer
	StaticPath      string
	RunID           int64
}

func Init() error {
	log.Printf("Starting...\n")

	Env.Ctx = context.Background()

	emailRecepients := os.Getenv("EMAIL_RECEPIENTS")
	if emailRecepients == "" {
		return fmt.Errorf("No recepients")
	}
	Env.Recepients = strings.Split(emailRecepients, ",")
	//fmt.Println(Env.Recepients)
	Env.ErrorRecepients = strings.Split(os.Getenv("ERROR_RECEPIENTS"), ",")
	Env.Sender = os.Getenv("EMAIL_SENDER")
	Env.SmtpHost = os.Getenv("EMAIL_HOST")
	Env.SmtpPass = os.Getenv("EMAIL_PASS")
	Env.StaticPath = os.Getenv("STATIC_PATH")

	d := gomail.NewDialer(Env.SmtpHost, 587, Env.Sender, Env.SmtpPass)
	Env.Smtp = d
	Env.Smtp.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	now := time.Now()
	Env.RunID = now.Unix()

	return nil
}
