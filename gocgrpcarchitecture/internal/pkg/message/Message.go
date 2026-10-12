package message

import (
	"bytes"
	"fmt"
	"gocgrpcarchitecture/internal/pkg/configs"
	"html/template"
	"log"
	"net/smtp"
)

var msg = configs.LoadConfig()

type Message struct {
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func Deliver(to []string, subject string) *Message {
	return &Message{
		to:      to,
		subject: subject,
	}
}

func (r *Message) sendMail() bool {
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%s", msg.Mail.Server, msg.Mail.Port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", msg.Mail.Email, msg.Mail.Password, msg.Mail.Server), msg.Mail.Email, r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

func (r *Message) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Message) EmailTemplate(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}
	if ok := r.sendMail(); ok {
		log.Printf("Email has been sent to %s\n", r.to)
	} else {
		log.Printf("Failed to send the email to %s\n", r.to)
	}
}
