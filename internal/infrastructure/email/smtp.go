package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"go.uber.org/zap"

	"github.com/agnathor/finances-go/internal/config"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/pkg/logger"
)

type SMTPSender struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSMTPSender(cfg config.EmailConfig) *SMTPSender {
	from := strings.TrimSpace(cfg.SMTPFrom)

	return &SMTPSender{
		host: cfg.SMTPHost,
		port: cfg.SMTPPort,
		user: cfg.SMTPUser,
		pass: cfg.SMTPPass,
		from: from,
	}
}

func (s *SMTPSender) IsConfigured() bool {
	return s.host != "" && s.port != "" && s.user != "" && s.pass != "" && s.from != ""
}

func (s *SMTPSender) Send(ctx context.Context, message domain.EmailMessage) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !s.IsConfigured() {
		logger.Get().Warn("smtp send skipped: sender is not configured",
			zap.Bool("smtp_host_configured", s.host != ""),
			zap.Bool("smtp_port_configured", s.port != ""),
			zap.Bool("smtp_user_configured", s.user != ""),
			zap.Bool("smtp_pass_configured", s.pass != ""),
			zap.Bool("smtp_from_configured", s.from != ""),
			zap.String("to", message.To),
			zap.String("subject", message.Subject),
		)
		return fmt.Errorf("smtp sender is not configured")
	}
	if strings.TrimSpace(message.To) == "" {
		logger.Get().Warn("smtp send skipped: recipient is empty",
			zap.String("from", s.from),
			zap.String("subject", message.Subject),
		)
		return fmt.Errorf("email recipient is required")
	}

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	body := buildMessage(s.from, message)
	log := logger.Get()

	if err := ctx.Err(); err != nil {
		return err
	}

	log.Info("smtp send started",
		zap.String("host", s.host),
		zap.String("port", s.port),
		zap.String("from", s.from),
		zap.String("to", message.To),
		zap.String("subject", message.Subject),
	)

	if err := smtp.SendMail(addr, auth, s.from, []string{message.To}, []byte(body)); err != nil {
		log.Warn("smtp send failed",
			zap.String("host", s.host),
			zap.String("port", s.port),
			zap.String("from", s.from),
			zap.String("to", message.To),
			zap.String("subject", message.Subject),
			zap.Error(err),
		)
		return err
	}

	log.Info("smtp send succeeded",
		zap.String("host", s.host),
		zap.String("port", s.port),
		zap.String("from", s.from),
		zap.String("to", message.To),
		zap.String("subject", message.Subject),
	)

	return nil
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
