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
	reminderRepo      domain.ReminderRepository
	emailSender       domain.EmailSender
	fallbackRecipient string
}

func NewService(reminderRepo domain.ReminderRepository, emailSender domain.EmailSender, fallbackRecipient string) Service {
	return &service{
		reminderRepo:      reminderRepo,
		emailSender:       emailSender,
		fallbackRecipient: strings.TrimSpace(fallbackRecipient),
	}
}

func (s *service) SendDue(ctx context.Context, today time.Time) (int, error) {
	log := logger.Get()
	log.Info("checking due reminder notifications",
		zap.String("timestamp", today.Format(time.RFC3339)),
		zap.String("date", today.Format("2006-01-02")),
		zap.String("time", today.Format("15:04")),
	)

	reminders, err := s.reminderRepo.GetDueForNotification(ctx, today)
	if err != nil {
		log.Warn("failed to fetch due reminder notifications", zap.Error(err))
		return 0, err
	}

	log.Info("due reminder notifications fetched", zap.Int("count", len(reminders)))

	sent := 0
	for _, reminder := range reminders {
		recipient := strings.TrimSpace(reminder.UserEmail)
		recipientSource := "user"
		if recipient == "" {
			recipient = s.fallbackRecipient
			recipientSource = "fallback"
		}
		if recipient == "" {
			log.Warn("skipping reminder notification: recipient email is missing",
				zap.String("reminder_id", reminder.ID),
				zap.String("user_id", reminder.UserID),
				zap.String("title", reminder.Title),
				zap.String("due_date", reminder.DueDate.Format("2006-01-02")),
				zap.String("reminder_time", reminder.ReminderTime),
			)
			continue
		}

		message := domain.EmailMessage{
			To:      recipient,
			Subject: fmt.Sprintf("Recordatorio: %s", reminder.Title),
			Body:    buildReminderBody(reminder),
		}

		log.Info("sending reminder notification email",
			zap.String("reminder_id", reminder.ID),
			zap.String("user_id", reminder.UserID),
			zap.String("to", recipient),
			zap.String("recipient_source", recipientSource),
			zap.String("subject", message.Subject),
			zap.String("title", reminder.Title),
			zap.String("due_date", reminder.DueDate.Format("2006-01-02")),
			zap.String("reminder_time", reminder.ReminderTime),
			zap.String("recurrence_type", string(reminder.RecurrenceType)),
		)

		if err := s.emailSender.Send(ctx, message); err != nil {
			log.Warn("failed to send reminder notification email",
				zap.String("reminder_id", reminder.ID),
				zap.String("user_id", reminder.UserID),
				zap.String("to", recipient),
				zap.String("subject", message.Subject),
				zap.Error(err),
			)
			return sent, err
		}

		log.Info("reminder notification email accepted by sender",
			zap.String("reminder_id", reminder.ID),
			zap.String("user_id", reminder.UserID),
			zap.String("to", recipient),
			zap.String("subject", message.Subject),
		)

		if err := s.reminderRepo.MarkNotificationSent(ctx, reminder.ID, today); err != nil {
			log.Warn("failed to mark reminder notification as sent",
				zap.String("reminder_id", reminder.ID),
				zap.String("user_id", reminder.UserID),
				zap.String("to", recipient),
				zap.Error(err),
			)
			continue
		}
		log.Info("reminder notification marked as sent",
			zap.String("reminder_id", reminder.ID),
			zap.String("user_id", reminder.UserID),
			zap.String("to", recipient),
			zap.String("sent_at", today.Format(time.RFC3339)),
		)
		sent++
	}

	log.Info("finished reminder notification run",
		zap.Int("due_count", len(reminders)),
		zap.Int("sent_count", sent),
	)

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
