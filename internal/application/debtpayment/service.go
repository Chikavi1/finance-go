package debtpayment

import (
	"context"
	"fmt"
	"time"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, debtID string, amount float64, paymentDate time.Time, notes *string) (*domain.DebtPayment, error)
	GetByDebtID(ctx context.Context, debtID string) ([]*domain.DebtPayment, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	debtPaymentRepo domain.DebtPaymentRepository
}

func NewService(debtPaymentRepo domain.DebtPaymentRepository) Service {
	return &service{debtPaymentRepo: debtPaymentRepo}
}

func (s *service) Create(ctx context.Context, debtID string, amount float64, paymentDate time.Time, notes *string) (*domain.DebtPayment, error) {
	payment := &domain.DebtPayment{
		DebtID:      debtID,
		Amount:      amount,
		PaymentDate: paymentDate,
		Notes:       notes,
	}

	if err := s.debtPaymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create debt payment: %w", err)
	}

	return payment, nil
}

func (s *service) GetByDebtID(ctx context.Context, debtID string) ([]*domain.DebtPayment, error) {
	return s.debtPaymentRepo.GetByDebtID(ctx, debtID)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.debtPaymentRepo.Delete(ctx, id)
}
