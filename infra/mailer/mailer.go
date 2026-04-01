package mailer

import (
	"CesarRodriguezPardo/template-go/config"
	logger "CesarRodriguezPardo/template-go/infra/logger"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strconv"

	"gopkg.in/mail.v2"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

const (
	DIR_GENERAL_URL  = "https://maps.app.goo.gl/YKHSrEKXFnBTRMzm7"
	DIR_GENERAL_TEXT = "El Belloto 3556, Estación Central."
)

var (
	email string
	pass  string
	host  string
	port  string
)

func InitMailer() {
	email = config.Cfg.Mailer.Dir
	pass = config.Cfg.Mailer.Pass
	host = config.Cfg.Mailer.Host
	port = config.Cfg.Mailer.Port
}

func NewRequest(to []string, subject string) *Request {
	return &Request{
		from:    "Sistema Gatitos y Usuarios",
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
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

func (r *Request) sendMail() error {
	body := "From: " + r.from + "\r\nTo: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%s", host, port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", email, pass, host), email, r.to, []byte(body)); err != nil {
		return err
	}
	return nil
}

func (r *Request) Send(templateName string, items interface{}) error {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		return err
	}
	if err := r.sendMail(); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *Request) SendMailSkipTLS(templateName string, items interface{}) error {

	err := r.parseTemplate(templateName, items)
	if err != nil {
		return err
	}

	m := mail.NewMessage()

	m.SetHeader("From", email)
	m.SetHeader("To", r.to[0]) //se pueden colocar mas mail separados por coma : m.SetHeader("To", "email1@email.com","email2@email.com")
	//m.SetAddressHeader("Cc", "manuel.manriquez.lopez@gmail.com", "Manuel")
	m.SetHeader("Subject", r.subject)
	m.SetBody("text/html", r.body)

	var puerto int
	puerto, err = strconv.Atoi(port)
	if err != nil {
		logger.Error("No se pudo mandar el email", err)
		return err
	}
	d := mail.NewDialer(host, puerto, email, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		logger.Error("No se pudo mandar el email", err)
		return err
	}

	return nil
}

func SendNewMail(to []string, subject string, templateName string, items interface{}) error {
	if os.Getenv("MAIL_DISABLED") == "true" {
		return nil
	}
	request := NewRequest(to, subject)
	err := request.Send(templateName, items)
	if err != nil {
		return err
	}
	return nil
}
