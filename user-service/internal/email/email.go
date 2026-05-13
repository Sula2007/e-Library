package email

import (
	"fmt"
	"net/smtp"
)

type Sender struct {
	host     string
	port     string
	username string
	password string
}

func NewSender(host, port, username, password string) *Sender {
	return &Sender{host: host, port: port, username: username, password: password}
}

func (s *Sender) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.username, to, subject, body)
	return smtp.SendMail(s.host+":"+s.port, auth, s.username, []string{to}, []byte(msg))
}