package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/application/account"
	"github.com/agnathor/finances-go/internal/application/attachment"
	"github.com/agnathor/finances-go/internal/application/auth"
	"github.com/agnathor/finances-go/internal/application/budget"
	"github.com/agnathor/finances-go/internal/application/category"
	"github.com/agnathor/finances-go/internal/application/dashboard"
	"github.com/agnathor/finances-go/internal/application/debt"
	debtpayment "github.com/agnathor/finances-go/internal/application/debtpayment"
	"github.com/agnathor/finances-go/internal/application/goal"
	goalcontribution "github.com/agnathor/finances-go/internal/application/goalcontribution"
	"github.com/agnathor/finances-go/internal/application/reminder"
	scheduledmovement "github.com/agnathor/finances-go/internal/application/scheduledmovement"
	"github.com/agnathor/finances-go/internal/application/setting"
	"github.com/agnathor/finances-go/internal/application/tag"
	"github.com/agnathor/finances-go/internal/application/transaction"
	"github.com/agnathor/finances-go/internal/application/user"
	"github.com/agnathor/finances-go/internal/config"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/infrastructure/database"
	"github.com/agnathor/finances-go/internal/infrastructure/storage"
	"github.com/agnathor/finances-go/internal/interfaces/api/handler"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/pkg/jwt"
)

type Dependencies struct {
	DB  *pgxpool.Pool
	Cfg config.Config
}

func NewRouter(deps Dependencies) *fiber.App {
	cfg := deps.Cfg

	app := fiber.New(fiber.Config{
		AppName:        cfg.App.Name,
		CaseSensitive:  true,
		StrictRouting:  false,
		ReadBufferSize: 8192,
		ErrorHandler:   errorHandler,
	})

	app.Use(requestid.New())
	app.Use(middleware.Recovery())
	app.Use(middleware.RequestLogger())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.App.CORSAllowedOrigins,
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Use(limiter.New(limiter.Config{
		Max: 100,
	}))

	jwtManager := jwt.NewManager(cfg.JWT)

	userRepo := database.NewUserRepository(deps.DB)
	refreshTokenRepo := database.NewRefreshTokenRepository(deps.DB)
	accountRepo := database.NewAccountRepository(deps.DB)
	categoryRepo := database.NewCategoryRepository(deps.DB)
	tagRepo := database.NewTagRepository(deps.DB)
	transactionRepo := database.NewTransactionRepository(deps.DB)
	attachmentRepo := database.NewAttachmentRepository(deps.DB)
	budgetRepo := database.NewBudgetRepository(deps.DB)
	goalRepo := database.NewGoalRepository(deps.DB)
	goalContributionRepo := database.NewGoalContributionRepository(deps.DB)
	debtRepo := database.NewDebtRepository(deps.DB)
	debtPaymentRepo := database.NewDebtPaymentRepository(deps.DB)
	reminderRepo := database.NewReminderRepository(deps.DB)
	scheduledMovementRepo := database.NewScheduledMovementRepository(deps.DB)
	settingRepo := database.NewSettingRepository(deps.DB)

	minioClient, err := storage.NewMinioClient(cfg.MinIO)
	var storageService domain.StorageService
	if err == nil {
		storageService = storage.NewStorageService(minioClient, cfg.MinIO.Bucket)
	}

	authService := auth.NewService(userRepo, refreshTokenRepo, jwtManager, cfg.JWT)
	userService := user.NewService(userRepo)
	accountService := account.NewService(accountRepo)
	categoryService := category.NewService(categoryRepo)
	tagService := tag.NewService(tagRepo)
	transactionService := transaction.NewService(transactionRepo)
	attachmentService := attachment.NewService(attachmentRepo, storageService)
	budgetService := budget.NewService(budgetRepo)
	goalService := goal.NewService(goalRepo)
	goalContributionService := goalcontribution.NewService(goalContributionRepo, goalRepo)
	debtServiceInstance := debt.NewService(debtRepo)
	debtPaymentService := debtpayment.NewService(debtPaymentRepo)
	reminderService := reminder.NewService(reminderRepo)
	scheduledMovementService := scheduledmovement.NewService(scheduledMovementRepo, transactionRepo)
	dashboardService := dashboard.NewService(accountRepo, transactionRepo, budgetRepo, debtRepo)
	settingService := setting.NewService(settingRepo)

	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	accountHandler := handler.NewAccountHandler(accountService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	attachmentHandler := handler.NewAttachmentHandler(attachmentService)
	budgetHandler := handler.NewBudgetHandler(budgetService)
	goalHandler := handler.NewGoalHandler(goalService, goalContributionService)
	debtHandler := handler.NewDebtHandler(debtServiceInstance, debtPaymentService)
	reminderHandler := handler.NewReminderHandler(reminderService)
	scheduledMovementHandler := handler.NewScheduledMovementHandler(scheduledMovementService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	settingHandler := handler.NewSettingHandler(settingService)

	api := app.Group("/api/v1")

	api.Get("/health", healthHandler.Check)

	authGroup := api.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.RefreshToken)
	authGroup.Post("/logout", middleware.AuthRequired(jwtManager), authHandler.Logout)

	userGroup := api.Group("/users", middleware.AuthRequired(jwtManager))
	userGroup.Get("/me", userHandler.GetProfile)
	userGroup.Put("/me", userHandler.UpdateProfile)
	userGroup.Put("/me/password", userHandler.ChangePassword)

	accountsGroup := api.Group("/accounts", middleware.AuthRequired(jwtManager))
	accountsGroup.Post("/", accountHandler.Create)
	accountsGroup.Get("/", accountHandler.GetAll)
	accountsGroup.Get("/:id", accountHandler.GetByID)
	accountsGroup.Put("/:id", accountHandler.Update)
	accountsGroup.Delete("/:id", accountHandler.Delete)

	categoriesGroup := api.Group("/categories", middleware.AuthRequired(jwtManager))
	categoriesGroup.Post("/", categoryHandler.Create)
	categoriesGroup.Get("/", categoryHandler.GetAll)
	categoriesGroup.Get("/:id", categoryHandler.GetByID)
	categoriesGroup.Put("/:id", categoryHandler.Update)
	categoriesGroup.Delete("/:id", categoryHandler.Delete)

	tagsGroup := api.Group("/tags", middleware.AuthRequired(jwtManager))
	tagsGroup.Post("/", tagHandler.Create)
	tagsGroup.Get("/", tagHandler.GetAll)
	tagsGroup.Get("/:id", tagHandler.GetByID)
	tagsGroup.Put("/:id", tagHandler.Update)
	tagsGroup.Delete("/:id", tagHandler.Delete)

	transactionsGroup := api.Group("/transactions", middleware.AuthRequired(jwtManager))
	transactionsGroup.Post("/", transactionHandler.Create)
	transactionsGroup.Get("/", transactionHandler.GetAll)
	transactionsGroup.Get("/:id", transactionHandler.GetByID)
	transactionsGroup.Put("/:id", transactionHandler.Update)
	transactionsGroup.Delete("/:id", transactionHandler.Delete)

	attachmentsGroup := api.Group("/attachments", middleware.AuthRequired(jwtManager))
	attachmentsGroup.Post("/upload", attachmentHandler.Upload)
	attachmentsGroup.Get("/:id", attachmentHandler.GetByID)
	attachmentsGroup.Get("/transaction/:transactionId", attachmentHandler.GetByTransactionID)
	attachmentsGroup.Delete("/:id", attachmentHandler.Delete)

	budgetsGroup := api.Group("/budgets", middleware.AuthRequired(jwtManager))
	budgetsGroup.Post("/", budgetHandler.Create)
	budgetsGroup.Get("/", budgetHandler.GetAll)
	budgetsGroup.Get("/:id", budgetHandler.GetByID)
	budgetsGroup.Put("/:id", budgetHandler.Update)
	budgetsGroup.Delete("/:id", budgetHandler.Delete)

	goalsGroup := api.Group("/goals", middleware.AuthRequired(jwtManager))
	goalsGroup.Post("/", goalHandler.Create)
	goalsGroup.Get("/", goalHandler.GetAll)
	goalsGroup.Get("/:id", goalHandler.GetByID)
	goalsGroup.Put("/:id", goalHandler.Update)
	goalsGroup.Delete("/:id", goalHandler.Delete)
	goalsGroup.Post("/:id/contributions", goalHandler.CreateContribution)
	goalsGroup.Get("/:id/contributions", goalHandler.GetContributions)
	goalsGroup.Delete("/:id/contributions/:contributionId", goalHandler.DeleteContribution)


	debtsGroup := api.Group("/debts", middleware.AuthRequired(jwtManager))
	debtsGroup.Post("/", debtHandler.Create)
	debtsGroup.Get("/", debtHandler.GetAll)
	debtsGroup.Get("/:id", debtHandler.GetByID)
	debtsGroup.Put("/:id", debtHandler.Update)
	debtsGroup.Delete("/:id", debtHandler.Delete)
	debtsGroup.Post("/:id/payments", debtHandler.CreatePayment)
	debtsGroup.Get("/:id/payments", debtHandler.GetPayments)
	debtsGroup.Delete("/:id/payments/:paymentId", debtHandler.DeletePayment)

	remindersGroup := api.Group("/reminders", middleware.AuthRequired(jwtManager))
	remindersGroup.Post("/", reminderHandler.Create)
	remindersGroup.Get("/", reminderHandler.GetAll)
	remindersGroup.Get("/:id", reminderHandler.GetByID)
	remindersGroup.Put("/:id", reminderHandler.Update)
	remindersGroup.Delete("/:id", reminderHandler.Delete)

	scheduledMovementsGroup := api.Group("/scheduled-movements", middleware.AuthRequired(jwtManager))
	scheduledMovementsGroup.Post("/", scheduledMovementHandler.Create)
	scheduledMovementsGroup.Get("/", scheduledMovementHandler.GetAll)
	scheduledMovementsGroup.Post("/generate-due", scheduledMovementHandler.GenerateDue)
	scheduledMovementsGroup.Get("/:id", scheduledMovementHandler.GetByID)
	scheduledMovementsGroup.Put("/:id", scheduledMovementHandler.Update)
	scheduledMovementsGroup.Delete("/:id", scheduledMovementHandler.Delete)

	dashboardGroup := api.Group("/dashboard", middleware.AuthRequired(jwtManager))
	dashboardGroup.Get("/", dashboardHandler.GetDashboard)

	settingsGroup := api.Group("/settings", middleware.AuthRequired(jwtManager))
	settingsGroup.Get("/", settingHandler.GetAll)
	settingsGroup.Put("/", settingHandler.Update)

	api.Get("/swagger/*", swagger.HandlerDefault)

	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}
