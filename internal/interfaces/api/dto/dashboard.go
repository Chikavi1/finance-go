package dto

type DashboardResponse struct {
	TotalBalance     float64               `json:"total_balance"`
	MonthlyIncome    float64               `json:"monthly_income"`
	MonthlyExpenses  float64               `json:"monthly_expenses"`
	RecentTxns       []TransactionResponse `json:"recent_transactions"`
	Budgets          []BudgetResponse      `json:"budgets"`
	ActiveDebtsTotal float64               `json:"active_debts_total"`
	ActiveDebtCount  int                   `json:"active_debt_count"`
}

type SettingResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateSettingsRequest struct {
	Settings map[string]string `json:"settings" validate:"required"`
}
