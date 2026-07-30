package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/agnathor/finances-go/internal/config"
	"github.com/agnathor/finances-go/internal/domain"
)

type SMTPSender struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSMTPSender(cfg config.EmailConfig) *SMTPSender {
	return &SMTPSender{
		host: cfg.SMTPHost,
		port: cfg.SMTPPort,
		user: cfg.SMTPUser,
		pass: cfg.SMTPPass,
		from: cfg.SMTPUser,
	}
}

func (s *SMTPSender) IsConfigured() bool {
	return s.host != "" && s.port != "" && s.user != "" && s.pass != ""
}

func (s *SMTPSender) Send(ctx context.Context, message domain.EmailMessage) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !s.IsConfigured() {
		return fmt.Errorf("smtp sender is not configured")
	}
	if strings.TrimSpace(message.To) == "" {
		return fmt.Errorf("email recipient is required")
	}

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	body := buildMessage(s.from, message)

	if err := ctx.Err(); err != nil {
		return err
	}
	return smtp.SendMail(addr, auth, s.from, []string{message.To}, []byte(body))
}

func buildMessage(from string, message domain.EmailMessage) string {
	headers := map[string]string{
		"From":         from,
		"To":           message.To,
		"Subject":      message.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=UTF-8",
	}

	var builder strings.Builder
	for key, value := range headers {
		builder.WriteString(key)
		builder.WriteString(": ")
		builder.WriteString(value)
		builder.WriteString("\r\n")
	}
	builder.WriteString("\r\n")
	builder.WriteString(message.Body)
	builder.WriteString("\r\n")
	return builder.String()
}
