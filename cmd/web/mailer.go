package main

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Wait        *sync.WaitGroup
	MailerChan  chan Message
	ErrorChan   chan error
	DoneChan    chan bool
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
	Template    string
}

func (m *Mail) SendMail(msg Message, errorChan chan error) {
	if msg.Template == "" {
		msg.Template = "mail"
	}
	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.BuildHTMLMessage(msg)
	if err != nil {
		errorChan <- err
	}
	plainMessage, err := m.BuildPlainTextMessage(msg)
	if err != nil {
		errorChan <- err
	}
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject).SetBody(mail.TextPlain, plainMessage).AddAlternative(mail.TextHTML, formattedMessage)
	if len(msg.Attachments) > 0 {
		for _, att := range msg.Attachments {
			email.AddAttachment(att)
		}
	}
	if err := email.Send(smtpClient); err != nil {
		errorChan <- err
	}
}

func (m *Mail) BuildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("../templates/%s.html.gohtml", msg.Template)
	t, err := template.New("email-html").ParseFiles(templateToRender)
	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}
	formattedMsg := tpl.String()
	formattedMsg, err = m.inLineCSS(formattedMsg)
	if err != nil {
		return "", err
	}
	return formattedMsg, nil
}

func (m *Mail) BuildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("../templates/%s.plain.gohtml", msg.Template)
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}
	formattedMsg := tpl.String()
	formattedMsg, err = m.inLineCSS(formattedMsg)
	if err != nil {
		return "", err
	}
	return formattedMsg, nil
}

func (m *Mail) inLineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}
	return html, nil
}

func (m *Mail) getEncryption(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
