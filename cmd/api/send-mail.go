package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/data"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
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
	if m.Template == "" {
		email.SetBody(mail.TextHTML, string(m.Content))
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./web/email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Printf("Failed to read template file: %v", err)
			return
		}
		mailTemplate := string(data)
		mailTemplate = strings.Replace(mailTemplate, "[%body%]", string(m.Content), 1)
		email.SetBody(mail.TextHTML, mailTemplate)
	}

	err = email.Send(smtpClient)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}
}
