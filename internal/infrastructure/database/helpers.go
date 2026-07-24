package database

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

func parseUUID(s string) (pgtype.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: u, Valid: true}, nil
}

func mustParseUUID(s string) pgtype.UUID {
	u := uuid.MustParse(s)
	return pgtype.UUID{Bytes: u, Valid: true}
}

func pgUUIDToString(u pgtype.UUID) string {
	return uuid.UUID(u.Bytes).String()
}

func toText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func fromText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func toTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func mapUser(u db.User) *domain.User {
	return &domain.User{
		ID:           pgUUIDToString(u.ID),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		AvatarURL:    fromText(u.AvatarUrl),
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}

func mapAccount(a db.Account) *domain.Account {
	return &domain.Account{
		ID:        pgUUIDToString(a.ID),
		UserID:    pgUUIDToString(a.UserID),
		Name:      a.Name,
		Type:      domain.AccountType(a.Type),
		Currency:  a.Currency,
		Balance:   a.Balance,
		Color:     a.Color,
		Icon:      a.Icon,
		Archived:  a.Archived,
		CreatedAt: a.CreatedAt.Time,
		UpdatedAt: a.UpdatedAt.Time,
	}
}

func mapCategory(c db.Category) *domain.Category {
	return &domain.Category{
		ID:        pgUUIDToString(c.ID),
		UserID:    pgUUIDToString(c.UserID),
		Name:      c.Name,
		Type:      domain.CategoryType(c.Type),
		Color:     c.Color,
		Icon:      c.Icon,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
}

func mapTag(t db.Tag) *domain.Tag {
	return &domain.Tag{
		ID:        pgUUIDToString(t.ID),
		UserID:    pgUUIDToString(t.UserID),
		Name:      t.Name,
		CreatedAt: t.CreatedAt.Time,
	}
}

func mapRefreshToken(t db.RefreshToken) *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:        pgUUIDToString(t.ID),
		UserID:    pgUUIDToString(t.UserID),
		TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt.Time,
		CreatedAt: t.CreatedAt.Time,
		Revoked:   t.Revoked,
	}
}

func mapBudget(b db.Budget) *domain.Budget {
	return &domain.Budget{
		ID:         pgUUIDToString(b.ID),
		UserID:     pgUUIDToString(b.UserID),
		CategoryID: pgUUIDToString(b.CategoryID),
		Amount:     b.Amount,
		Spent:      b.Spent,
		Month:      b.Month,
		Year:       b.Year,
		CreatedAt:  b.CreatedAt.Time,
		UpdatedAt:  b.UpdatedAt.Time,
	}
}

func mapGoal(g db.Goal) *domain.Goal {
	var targetDate *time.Time
	if g.TargetDate.Valid {
		targetDate = &g.TargetDate.Time
	}

	return &domain.Goal{
		ID:            pgUUIDToString(g.ID),
		UserID:        pgUUIDToString(g.UserID),
		Name:          g.Name,
		TargetAmount:  g.TargetAmount,
		CurrentAmount: g.CurrentAmount,
		TargetDate:    targetDate,
		Icon:          g.Icon,
		Color:         g.Color,
		CreatedAt:     g.CreatedAt.Time,
		UpdatedAt:     g.UpdatedAt.Time,
	}
}

func mapDebt(d db.Debt) *domain.Debt {
	var dueDate *time.Time
	if d.DueDate.Valid {
		dueDate = &d.DueDate.Time
	}

	return &domain.Debt{
		ID:              pgUUIDToString(d.ID),
		UserID:          pgUUIDToString(d.UserID),
		Name:            d.Name,
		TotalAmount:     d.TotalAmount,
		RemainingAmount: d.RemainingAmount,
		InterestRate:    d.InterestRate,
		DueDate:         dueDate,
		Status:          domain.DebtStatus(d.Status),
		Notes:           fromText(d.Notes),
		CreatedAt:       d.CreatedAt.Time,
		UpdatedAt:       d.UpdatedAt.Time,
	}
}

func mapDebtPayment(p db.DebtPayment) *domain.DebtPayment {
	return &domain.DebtPayment{
		ID:          pgUUIDToString(p.ID),
		DebtID:      pgUUIDToString(p.DebtID),
		Amount:      p.Amount,
		PaymentDate: p.PaymentDate.Time,
		Notes:       fromText(p.Notes),
		CreatedAt:   p.CreatedAt.Time,
	}
}

func toNullableDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}
