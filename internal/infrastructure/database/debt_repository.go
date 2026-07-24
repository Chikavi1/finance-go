package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type debtRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewDebtRepository(pool *pgxpool.Pool) domain.DebtRepository {
	return &debtRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *debtRepository) Create(ctx context.Context, debt *domain.Debt) error {
	userUUID, err := parseUUID(debt.UserID)
	if err != nil {
		return domain.ErrNotFound
	}

	created, err := r.query.CreateDebt(ctx, db.CreateDebtParams{
		UserID:          userUUID,
		Name:            debt.Name,
		TotalAmount:     debt.TotalAmount,
		RemainingAmount: debt.RemainingAmount,
		InterestRate:    debt.InterestRate,
		DueDate:         toNullableDate(debt.DueDate),
		Status:          string(debt.Status),
		Notes:           toText(debt.Notes),
	})
	if err != nil {
		return err
	}

	debt.ID = pgUUIDToString(created.ID)
	debt.CreatedAt = created.CreatedAt.Time
	debt.UpdatedAt = created.UpdatedAt.Time
	return nil
}

func (r *debtRepository) GetByID(ctx context.Context, id string) (*domain.Debt, error) {
	debtUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	debt, err := r.query.GetDebtByID(ctx, debtUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapDebt(debt), nil
}

func (r *debtRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Debt, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	debts, err := r.query.GetDebtsByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Debt, len(debts))
	for i, d := range debts {
		result[i] = mapDebt(d)
	}

	return result, nil
}

func (r *debtRepository) Update(ctx context.Context, debt *domain.Debt) error {
	debtUUID, err := parseUUID(debt.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	updated, err := r.query.UpdateDebt(ctx, db.UpdateDebtParams{
		ID:              debtUUID,
		Name:            debt.Name,
		TotalAmount:     debt.TotalAmount,
		RemainingAmount: debt.RemainingAmount,
		InterestRate:    debt.InterestRate,
		DueDate:         toNullableDate(debt.DueDate),
		Status:          string(debt.Status),
		Notes:           toText(debt.Notes),
	})
	if err != nil {
		if isNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}

	debt.UpdatedAt = updated.UpdatedAt.Time
	return nil
}

func (r *debtRepository) Delete(ctx context.Context, id string) error {
	debtUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteDebt(ctx, debtUUID)
}

type debtPaymentRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewDebtPaymentRepository(pool *pgxpool.Pool) domain.DebtPaymentRepository {
	return &debtPaymentRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *debtPaymentRepository) Create(ctx context.Context, payment *domain.DebtPayment) error {
	debtUUID, err := parseUUID(payment.DebtID)
	if err != nil {
		return domain.ErrNotFound
	}

	created, err := r.query.CreateDebtPayment(ctx, db.CreateDebtPaymentParams{
		DebtID:      debtUUID,
		Amount:      payment.Amount,
		PaymentDate: pgtype.Date{Time: payment.PaymentDate, Valid: true},
		Notes:       toText(payment.Notes),
	})
	if err != nil {
		return err
	}

	payment.ID = pgUUIDToString(created.ID)
	payment.CreatedAt = created.CreatedAt.Time
	return nil
}

func (r *debtPaymentRepository) GetByDebtID(ctx context.Context, debtID string) ([]*domain.DebtPayment, error) {
	debtUUID, err := parseUUID(debtID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	payments, err := r.query.GetDebtPaymentsByDebtID(ctx, debtUUID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.DebtPayment, len(payments))
	for i, p := range payments {
		result[i] = mapDebtPayment(p)
	}

	return result, nil
}

func (r *debtPaymentRepository) Delete(ctx context.Context, id string) error {
	paymentUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteDebtPayment(ctx, paymentUUID)
}
