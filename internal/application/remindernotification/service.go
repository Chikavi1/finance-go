package remindernotification

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/pkg/logger"
)

type Service interface {
	SendDue(ctx context.Context, today time.Time) (int, error)
}

type service struct {
	reminderRepo domain.ReminderRepository
	emailSender  domain.EmailSender
	recipient    string
}

func NewService(reminderRepo domain.ReminderRepository, emailSender domain.EmailSender, recipient string) Service {
	return &service{
		reminderRepo: reminderRepo,
		emailSender:  emailSender,
		recipient:    strings.TrimSpace(recipient),
	}
}

func (s *service) SendDue(ctx context.Context, today time.Time) (int, error) {
	if s.recipient == "" {
		return 0, fmt.Errorf("report notification email is not configured")
	}

	reminders, err := s.reminderRepo.GetDueForNotification(ctx, today)
	if err != nil {
		return 0, err
	}

	sent := 0
	for _, reminder := range reminders {
		message := domain.EmailMessage{
			To:      s.recipient,
			Subject: fmt.Sprintf("Recordatorio: %s", reminder.Title),
			Body:    buildReminderBody(reminder),
		}
		if err := s.emailSender.Send(ctx, message); err != nil {
			return sent, err
		}

		if err := s.reminderRepo.MarkNotificationSent(ctx, reminder.ID, today); err != nil {
			logger.Get().Warn("failed to mark reminder notification as sent",
				zap.String("reminder_id", reminder.ID),
				zap.Error(err),
			)
			continue
		}
		sent++
	}

	return sent, nil
}

func buildReminderBody(reminder *domain.Reminder) string {
	var builder strings.Builder
	builder.WriteString("Tienes un recordatorio pendiente.\n\n")
	builder.WriteString("Concepto: ")
	builder.WriteString(reminder.Title)
	builder.WriteString("\n")
	builder.WriteString("Fecha: ")
	builder.WriteString(reminder.DueDate.Format("2006-01-02"))
	builder.WriteString("\n")
	builder.WriteString("Hora: ")
	builder.WriteString(reminder.ReminderTime)
	builder.WriteString("\n")
	if reminder.RecurrenceType == domain.ReminderRecurrenceMonthly && reminder.DayOfMonth != nil {
		builder.WriteString("Repetición: cada mes, día ")
		builder.WriteString(fmt.Sprintf("%d", *reminder.DayOfMonth))
		builder.WriteString("\n")
	}

	if reminder.Amount != nil {
		builder.WriteString("Monto: $")
		builder.WriteString(fmt.Sprintf("%.2f", *reminder.Amount))
		builder.WriteString("\n")
	}
	if reminder.Notes != nil && strings.TrimSpace(*reminder.Notes) != "" {
		builder.WriteString("\nNotas:\n")
		builder.WriteString(*reminder.Notes)
		builder.WriteString("\n")
	}

	builder.WriteString("\nFinai")
	return builder.String()
}
