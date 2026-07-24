package dashboard

import (
	"context"
	"time"

	"github.com/agnathor/finances-go/internal/domain"
)

type DashboardData struct {
	TotalBalance    float64              `json:"total_balance"`
	MonthlyIncome   float64              `json:"monthly_income"`
	MonthlyExpenses float64              `json:"monthly_expenses"`
	RecentTxns      []*domain.Transaction `json:"recent_transactions"`
	Budgets         []*domain.Budget      `json:"budgets"`
	ActiveDebtsTotal float64             `json:"active_debts_total"`
	ActiveDebtCount int                  `json:"active_debt_count"`
}

type Service interface {
	GetDashboard(ctx context.Context, userID string) (*DashboardData, error)
}

type service struct {
	accountRepo     domain.AccountRepository
	transactionRepo domain.TransactionRepository
	budgetRepo      domain.BudgetRepository
	debtRepo        domain.DebtRepository
}

func NewService(
	accountRepo domain.AccountRepository,
	transactionRepo domain.TransactionRepository,
	budgetRepo domain.BudgetRepository,
	debtRepo domain.DebtRepository,
) Service {
	return &service{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		budgetRepo:      budgetRepo,
		debtRepo:        debtRepo,
	}
}

func (s *service) GetDashboard(ctx context.Context, userID string) (*DashboardData, error) {
	accounts, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalBalance := 0.0
	for _, a := range accounts {
		totalBalance += a.Balance
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	startOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	incomeFilter := domain.TransactionFilter{
		Type:      incomeType(),
		StartDate: &startOfMonth,
		EndDate:   &endOfMonth,
	}

	expenseFilter := domain.TransactionFilter{
		Type:      expenseType(),
		StartDate: &startOfMonth,
		EndDate:   &endOfMonth,
	}

	allMonthFilter := domain.TransactionFilter{
		StartDate: &startOfMonth,
		EndDate:   &endOfMonth,
	}

	incomeTxns, _ := s.transactionRepo.GetByUserID(ctx, userID, incomeFilter)
	expenseTxns, _ := s.transactionRepo.GetByUserID(ctx, userID, expenseFilter)

	monthlyIncome := 0.0
	for _, t := range incomeTxns {
		monthlyIncome += t.Amount
	}

	monthlyExpenses := 0.0
	for _, t := range expenseTxns {
		monthlyExpenses += t.Amount
	}

	recentTxns, _ := s.transactionRepo.GetByUserID(ctx, userID, allMonthFilter)
	if len(recentTxns) > 5 {
		recentTxns = recentTxns[:5]
	}

	budgets, _ := s.budgetRepo.GetByMonthYear(ctx, userID, int32(currentMonth), int32(currentYear))

	allDebts, _ := s.debtRepo.GetByUserID(ctx, userID)
	activeDebtsTotal := 0.0
	activeDebtCount := 0
	for _, d := range allDebts {
		if d.Status == domain.DebtStatusActive || d.Status == domain.DebtStatusOverdue {
			activeDebtsTotal += d.RemainingAmount
			activeDebtCount++
		}
	}

	return &DashboardData{
		TotalBalance:     totalBalance,
		MonthlyIncome:    monthlyIncome,
		MonthlyExpenses:  monthlyExpenses,
		RecentTxns:       recentTxns,
		Budgets:          budgets,
		ActiveDebtsTotal: activeDebtsTotal,
		ActiveDebtCount:  activeDebtCount,
	}, nil
}

func incomeType() *domain.TransactionType {
	t := domain.TransactionTypeIncome
	return &t
}

func expenseType() *domain.TransactionType {
	t := domain.TransactionTypeExpense
	return &t
}
