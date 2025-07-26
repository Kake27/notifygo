package service

import (
	"fmt"
	"net/smtp"
	"notification-service/config"
)

type EmailService struct {
	cfg config.SMTPConfig
}

func NewEmailService() *EmailService {
	return &EmailService{
		cfg: config.GetSMTPConfig(),
	}
}

func (e *EmailService) Send(to, subject, message string) error {
	addr := e.cfg.Address()

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, message))
	err := smtp.SendMail(addr, nil, e.cfg.Sender, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)  
	}

	return nil
}