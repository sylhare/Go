package adapters

import "fmt"

type ConsoleWelcomeMailer struct{}

func NewConsoleWelcomeMailer() *ConsoleWelcomeMailer {
    return &ConsoleWelcomeMailer{}
}

func (m *ConsoleWelcomeMailer) SendWelcome(toEmail, name string) error {
    fmt.Printf("[CONSOLE MAILER] Sending welcome email to %s <%s>\n", name, toEmail)
    return nil
}
