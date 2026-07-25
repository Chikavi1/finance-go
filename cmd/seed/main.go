package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/agnathor/finances-go/internal/config"
	database "github.com/agnathor/finances-go/internal/infrastructure/database"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
	"github.com/agnathor/finances-go/pkg/password"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	q := db.New(pool)
	ctx := context.Background()

	fmt.Println("Seeding database...")

	userID := seedUser(ctx, q)
	fmt.Printf("  User created: id=%s email=chikavi10@gmail.com password=12345678A\n", userID)

	accIDs := seedAccounts(ctx, q, userID)
	fmt.Println("  Accounts created: 4")

	catIncomeIDs, catExpenseIDs := seedCategories(ctx, q, userID)
	fmt.Println("  Categories created: 9")

	tagIDs := seedTags(ctx, q, userID)
	fmt.Println("  Tags created: 5")

	seedTransactions(ctx, q, userID, accIDs, catIncomeIDs, catExpenseIDs, tagIDs)
	fmt.Println("  Transactions created: 15")

	seedBudgets(ctx, q, userID, catExpenseIDs)
	fmt.Println("  Budgets created: 6")

	seedGoals(ctx, q, userID)
	fmt.Println("  Goals created: 3")

	seedDebts(ctx, q, userID)
	fmt.Println("  Debts created: 1")

	seedSettings(ctx, q, userID)
	fmt.Println("  Settings created: 3")

	fmt.Println("Done!")
}

func parseUUID(s string) pgtype.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		log.Fatalf("invalid UUID %s: %v", s, err)
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}

func toText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func toNullUUID() pgtype.UUID {
	return pgtype.UUID{Valid: false}
}

func toDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func seedUser(ctx context.Context, q *db.Queries) string {
	hashed, err := password.Hash("12345678A")
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	u, err := q.CreateUser(ctx, db.CreateUserParams{
		Email:        "chikavi10@gmail.com",
		PasswordHash: hashed,
		Name:         "Usuario Demo",
		AvatarUrl:    toText(nil),
	})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	return uuid.UUID(u.ID.Bytes).String()
}

func seedAccounts(ctx context.Context, q *db.Queries, userIDStr string) map[string]string {
	uid := parseUUID(userIDStr)
	accounts := []struct {
		key      string
		name     string
		accType  string
		currency string
		balance  float64
		color    string
		icon     string
	}{
		{"cash", "Efectivo", "cash", "MXN", 10000, "#22c55e", "wallet"},
		{"bank", "Cuenta de Nómina", "bank", "MXN", 25000, "#3b82f6", "bank"},
		{"credit", "Tarjeta de Crédito", "credit_card", "MXN", 0, "#ef4444", "credit-card"},
		{"savings", "Ahorros", "savings", "MXN", 50000, "#8b5cf6", "piggy-bank"},
	}
	ids := make(map[string]string, len(accounts))
	for _, a := range accounts {
		created, err := q.CreateAccount(ctx, db.CreateAccountParams{
			UserID:   uid,
			Name:     a.name,
			Type:     a.accType,
			Currency: a.currency,
			Balance:  a.balance,
			Color:    a.color,
			Icon:     a.icon,
		})
		if err != nil {
			log.Fatalf("failed to create account %s: %v", a.name, err)
		}
		ids[a.key] = uuid.UUID(created.ID.Bytes).String()
	}
	return ids
}

func seedCategories(ctx context.Context, q *db.Queries, userIDStr string) (incomeIDs, expenseIDs map[string]string) {
	uid := parseUUID(userIDStr)

	income := []struct {
		key   string
		name  string
		color string
		icon  string
	}{
		{"salary", "Salario", "#22c55e", "briefcase"},
		{"freelance", "Freelance", "#06b6d4", "laptop"},
		{"investments", "Inversiones", "#8b5cf6", "trending-up"},
	}
	expense := []struct {
		key   string
		name  string
		color string
		icon  string
	}{
		{"food", "Alimentación", "#ef4444", "food"},
		{"transport", "Transporte", "#f97316", "car"},
		{"housing", "Vivienda", "#eab308", "home"},
		{"entertainment", "Entretenimiento", "#ec4899", "gamepad-2"},
		{"health", "Salud", "#06b6d4", "heart-pulse"},
		{"education", "Educación", "#8b5cf6", "book"},
	}

	incomeIDs = make(map[string]string, len(income))
	expenseIDs = make(map[string]string, len(expense))

	for _, c := range income {
		created, err := q.CreateCategory(ctx, db.CreateCategoryParams{
			UserID: uid, Name: c.name, Type: "income", Color: c.color, Icon: c.icon,
		})
		if err != nil {
			log.Fatalf("failed to create income category %s: %v", c.name, err)
		}
		incomeIDs[c.key] = uuid.UUID(created.ID.Bytes).String()
	}
	for _, c := range expense {
		created, err := q.CreateCategory(ctx, db.CreateCategoryParams{
			UserID: uid, Name: c.name, Type: "expense", Color: c.color, Icon: c.icon,
		})
		if err != nil {
			log.Fatalf("failed to create expense category %s: %v", c.name, err)
		}
		expenseIDs[c.key] = uuid.UUID(created.ID.Bytes).String()
	}
	return incomeIDs, expenseIDs
}

func seedTags(ctx context.Context, q *db.Queries, userIDStr string) map[string]string {
	uid := parseUUID(userIDStr)
	tagNames := []string{"esencial", "suscripcion", "emergencia", "ocio", "anual"}
	ids := make(map[string]string, len(tagNames))
	for _, name := range tagNames {
		created, err := q.CreateTag(ctx, db.CreateTagParams{UserID: uid, Name: name})
		if err != nil {
			log.Fatalf("failed to create tag %s: %v", name, err)
		}
		ids[name] = uuid.UUID(created.ID.Bytes).String()
	}
	return ids
}

func seedTransactions(ctx context.Context, q *db.Queries, userIDStr string, accIDs, catIncomeIDs, catExpenseIDs, tagIDs map[string]string) {
	uid := parseUUID(userIDStr)
	now := time.Now()

	type txDef struct {
		date        time.Time
		desc        string
		txType      string
		accountKey  string
		categoryKey string
		amount      float64
		tagKeys     []string
	}

	txs := []txDef{
		{now.AddDate(0, 0, -1), "Salario mensual", "income", "bank", "salary", 15000, []string{"esencial"}},
		{now.AddDate(0, 0, -2), "Supermercado", "expense", "credit", "food", 850, []string{"esencial"}},
		{now.AddDate(0, 0, -2), "Gasolina", "expense", "credit", "transport", 500, []string{"esencial"}},
		{now.AddDate(0, 0, -3), "Netflix", "expense", "credit", "entertainment", 239, []string{"suscripcion", "ocio"}},
		{now.AddDate(0, 0, -3), "Renta departamento", "expense", "bank", "housing", 4500, []string{"esencial"}},
		{now.AddDate(0, 0, -5), "Curso online", "expense", "credit", "education", 1200, []string{"anual"}},
		{now.AddDate(0, 0, -7), "Freelance diseño web", "income", "bank", "freelance", 5000, []string{}},
		{now.AddDate(0, 0, -8), "Cena restaurante", "expense", "credit", "food", 650, []string{"ocio"}},
		{now.AddDate(0, 0, -10), "Uber", "expense", "credit", "transport", 180, []string{}},
		{now.AddDate(0, 0, -12), "Spotify", "expense", "credit", "entertainment", 129, []string{"suscripcion", "ocio"}},
		{now.AddDate(0, 0, -14), "Consulta médica", "expense", "cash", "health", 800, []string{"esencial"}},
		{now.AddDate(0, 0, -15), "Dividendos inversiones", "income", "savings", "investments", 1200, []string{}},
		{now.AddDate(0, 0, -20), "Libros", "expense", "credit", "education", 450, []string{"ocio"}},
		{now.AddDate(0, 0, -25), "Agua", "expense", "bank", "housing", 350, []string{"esencial"}},
		{now.AddDate(0, 0, -28), "Luz", "expense", "bank", "housing", 420, []string{"esencial"}},
	}

	for _, t := range txs {
		var catID pgtype.UUID
		if t.txType == "income" {
			if id, ok := catIncomeIDs[t.categoryKey]; ok {
				catID = parseUUID(id)
			}
		} else {
			if id, ok := catExpenseIDs[t.categoryKey]; ok {
				catID = parseUUID(id)
			}
		}

		created, err := q.CreateTransaction(ctx, db.CreateTransactionParams{
			UserID:      uid,
			AccountID:   parseUUID(accIDs[t.accountKey]),
			ToAccountID: toNullUUID(),
			CategoryID:  catID,
			Type:        t.txType,
			Amount:      t.amount,
			Description: t.desc,
			Notes:       toText(nil),
			Date:        toDate(t.date),
		})
		if err != nil {
			log.Fatalf("failed to create transaction %s: %v", t.desc, err)
		}

		for _, tagKey := range t.tagKeys {
			if tagID, ok := tagIDs[tagKey]; ok {
				if err := q.CreateTransactionTag(ctx, db.CreateTransactionTagParams{
					TransactionID: created.ID,
					TagID:         parseUUID(tagID),
				}); err != nil {
					log.Fatalf("failed to add tag to %s: %v", t.desc, err)
				}
			}
		}
	}
}

func seedBudgets(ctx context.Context, q *db.Queries, userIDStr string, catExpenseIDs map[string]string) {
	uid := parseUUID(userIDStr)
	now := time.Now()
	month := int32(now.Month())
	year := int32(now.Year())

	budgets := []struct {
		categoryKey string
		amount      float64
	}{
		{"food", 4000}, {"transport", 2500}, {"housing", 7000},
		{"entertainment", 1500}, {"health", 2000}, {"education", 2000},
	}
	for _, b := range budgets {
		_, err := q.CreateBudget(ctx, db.CreateBudgetParams{
			UserID:     uid,
			CategoryID: parseUUID(catExpenseIDs[b.categoryKey]),
			Amount:     b.amount,
			Spent:      b.amount * 0.6,
			Month:      month,
			Year:       year,
		})
		if err != nil {
			log.Fatalf("failed to create budget for %s: %v", b.categoryKey, err)
		}
	}
}

func seedGoals(ctx context.Context, q *db.Queries, userIDStr string) {
	uid := parseUUID(userIDStr)
	goals := []struct {
		name          string
		targetAmount  float64
		currentAmount float64
		targetDate    time.Time
		icon          string
		color         string
	}{
		{"Viaje a Europa", 200000, 50000, time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC), "plane", "#3b82f6"},
		{"Fondo de Emergencia", 100000, 30000, time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC), "shield", "#22c55e"},
		{"Auto Nuevo", 500000, 100000, time.Date(2028, 12, 31, 0, 0, 0, 0, time.UTC), "car", "#f97316"},
	}
	for _, g := range goals {
		_, err := q.CreateGoal(ctx, db.CreateGoalParams{
			UserID:        uid,
			Name:          g.name,
			TargetAmount:  g.targetAmount,
			CurrentAmount: g.currentAmount,
			TargetDate:    toDate(g.targetDate),
			Icon:          g.icon,
			Color:         g.color,
		})
		if err != nil {
			log.Fatalf("failed to create goal %s: %v", g.name, err)
		}
	}
}

func seedDebts(ctx context.Context, q *db.Queries, userIDStr string) {
	uid := parseUUID(userIDStr)
	notes := "Préstamo para renovación de oficina"

	d, err := q.CreateDebt(ctx, db.CreateDebtParams{
		UserID:          uid,
		Name:            "Préstamo Personal",
		TotalAmount:     50000,
		RemainingAmount: 30000,
		InterestRate:    12.0,
		DueDate:         toDate(time.Date(2027, 6, 15, 0, 0, 0, 0, time.UTC)),
		Status:          "active",
		Notes:           toText(&notes),
	})
	if err != nil {
		log.Fatalf("failed to create debt: %v", err)
	}

	payments := []struct {
		amount      float64
		paymentDate time.Time
		notes       string
	}{
		{5000, time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC), "Pago mensual abril"},
		{5000, time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC), "Pago mensual mayo"},
		{5000, time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC), "Pago mensual junio"},
		{5000, time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC), "Pago mensual julio"},
	}
	for _, p := range payments {
		notes := p.notes
		_, err := q.CreateDebtPayment(ctx, db.CreateDebtPaymentParams{
			DebtID:      d.ID,
			Amount:      p.amount,
			PaymentDate: toDate(p.paymentDate),
			Notes:       toText(&notes),
		})
		if err != nil {
			log.Fatalf("failed to create debt payment: %v", err)
		}
	}
}

func seedSettings(ctx context.Context, q *db.Queries, userIDStr string) {
	uid := parseUUID(userIDStr)
	settings := []struct{ key, value string }{
		{"currency", "MXN"},
		{"theme", "system"},
		{"language", "es"},
	}
	for _, s := range settings {
		_, err := q.UpsertSetting(ctx, db.UpsertSettingParams{
			UserID: uid, Key: s.key, Value: s.value,
		})
		if err != nil {
			log.Fatalf("failed to upsert setting %s: %v", s.key, err)
		}
	}
}
