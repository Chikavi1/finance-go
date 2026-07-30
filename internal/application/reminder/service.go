package reminder

import (
	"context"
	"fmt"
	"time"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID, title string, amount *float64, dueDate time.Time, reminderTime string, recurrenceType domain.ReminderRecurrenceType, dayOfMonth *int, status domain.ReminderStatus, relatedType, relatedID, notes *string) (*domain.Reminder, error)
	GetByID(ctx context.Context, userID, id string) (*domain.Reminder, error)
	GetAll(ctx context.Context, userID string, includeDone bool) ([]*domain.Reminder, error)
	Update(ctx context.Context, userID, id, title string, amount *float64, dueDate time.Time, reminderTime string, recurrenceType domain.ReminderRecurrenceType, dayOfMonth *int, status domain.ReminderStatus, relatedType, relatedID, notes *string) (*domain.Reminder, error)
	Delete(ctx context.Context, userID, id string) error
}

type service struct {
	reminderRepo domain.ReminderRepository
}

func NewService(reminderRepo domain.ReminderRepository) Service {
	return &service{reminderRepo: reminderRepo}
}

func (s *service) Create(ctx context.Context, userID, title string, amount *float64, dueDate time.Time, reminderTime string, recurrenceType domain.ReminderRecurrenceType, dayOfMonth *int, status domain.ReminderStatus, relatedType, relatedID, notes *string) (*domain.Reminder, error) {
	if status == "" {
		status = domain.ReminderStatusPending
	}
	if reminderTime == "" {
		reminderTime = "09:00"
	}
	if recurrenceType == "" {
		recurrenceType = domain.ReminderRecurrenceOnce
	}
	if dayOfMonth == nil {
		day := dueDate.Day()
		dayOfMonth = &day
	}

	reminder := &domain.Reminder{
		UserID:      userID,
		Title:       title,
		Amount:      amount,
		DueDate:     dueDate,
		ReminderTime: reminderTime,
		RecurrenceType: recurrenceType,
		DayOfMonth: dayOfMonth,
		Status:      status,
		RelatedType: relatedType,
		RelatedID:   relatedID,
		Notes:       notes,
	}

	if err := s.reminderRepo.Create(ctx, reminder); err != nil {
		return nil, fmt.Errorf("failed to create reminder: %w", err)
	}

	return reminder, nil
}

func (s *service) GetByID(ctx context.Context, userID, id string) (*domain.Reminder, error) {
	reminder, err := s.reminderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reminder.UserID != userID {
		return nil, domain.ErrNotFound
	}
	return reminder, nil
}

func (s *service) GetAll(ctx context.Context, userID string, includeDone bool) ([]*domain.Reminder, error) {
	return s.reminderRepo.GetByUserID(ctx, userID, includeDone)
}

func (s *service) Update(ctx context.Context, userID, id, title string, amount *float64, dueDate time.Time, reminderTime string, recurrenceType domain.ReminderRecurrenceType, dayOfMonth *int, status domain.ReminderStatus, relatedType, relatedID, notes *string) (*domain.Reminder, error) {
	reminder, err := s.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if status == "" {
		status = domain.ReminderStatusPending
	}
	if reminderTime == "" {
		reminderTime = "09:00"
	}
	if recurrenceType == "" {
		recurrenceType = domain.ReminderRecurrenceOnce
	}
	if dayOfMonth == nil {
		day := dueDate.Day()
		dayOfMonth = &day
	}

	reminder.Title = title
	reminder.Amount = amount
	reminder.DueDate = dueDate
	reminder.ReminderTime = reminderTime
	reminder.RecurrenceType = recurrenceType
	reminder.DayOfMonth = dayOfMonth
	reminder.Status = status
	reminder.RelatedType = relatedType
	reminder.RelatedID = relatedID
	reminder.Notes = notes

	if err := s.reminderRepo.Update(ctx, reminder); err != nil {
		return nil, fmt.Errorf("failed to update reminder: %w", err)
	}

	return reminder, nil
}

func (s *service) Delete(ctx context.Context, userID, id string) error {
	if _, err := s.GetByID(ctx, userID, id); err != nil {
		return err
	}
	return s.reminderRepo.Delete(ctx, id)
}
