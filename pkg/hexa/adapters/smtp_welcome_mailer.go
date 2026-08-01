package adapters

import (
    "fmt"
    "net/smtp"
)

type SmtpWelcomeMailer struct {
    addr string
}

func NewSmtpWelcomeMailer(addr string) *SmtpWelcomeMailer {
    return &SmtpWelcomeMailer{addr: addr}
}

func (m *SmtpWelcomeMailer) SendWelcome(toEmail, name string) error {
    body := fmt.Sprintf("Subject: Welcome!\r\n\r\nHello %s, welcome aboard.", name)
    return smtp.SendMail(m.addr, nil, "noreply@app.com", []string{toEmail}, []byte(body))
}
