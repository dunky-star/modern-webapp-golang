package main

import (
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/data"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listerForMail() {
	app.InfoLog.Println("Email listener started - ready to send emails")
	go func() {
		for {
			msg := <-app.MailChan
			sendMail(msg)
			app.InfoLog.Println("Mail sent to", msg.To)
		}
	}()
}

func sendMail(m data.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		app.ErrorLog.Printf("Failed to connect to SMTP server: %v", err)
		return
	}
	defer smtpClient.Close()

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	// Convert template.HTML to string for email library (template.HTML is a type-safe alias for trusted HTML)
	email.SetBody(mail.TextHTML, string(m.Content))

	err = email.Send(smtpClient)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}
}
