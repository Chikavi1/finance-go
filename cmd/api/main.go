package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	_ "github.com/agnathor/finances-go/docs"
	reminderNotification "github.com/agnathor/finances-go/internal/application/remindernotification"
	"github.com/agnathor/finances-go/internal/config"
	"github.com/agnathor/finances-go/internal/infrastructure/cache"
	"github.com/agnathor/finances-go/internal/infrastructure/database"
	emailInfra "github.com/agnathor/finances-go/internal/infrastructure/email"
	api "github.com/agnathor/finances-go/internal/interfaces/api"
	"github.com/agnathor/finances-go/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger.Init(cfg.App.LogLevel)
	defer logger.Sync()

	log := logger.Get()
	log.Info("starting server", zap.String("env", cfg.App.Env))

	db, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatal("failed to run migrations", zap.Error(err))
	}
	database.EnsureReminderSchema(context.Background(), db, log)

	rdb, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Warn("redis not available, continuing without cache", zap.Error(err))
	}
	if rdb != nil {
		defer rdb.Close()
	}

	reminderRepo := database.NewReminderRepository(db)
	emailSender := emailInfra.NewSMTPSender(cfg.Email)
	workerCtx, stopReminderWorker := context.WithCancel(context.Background())
	defer stopReminderWorker()
	notificationService := reminderNotification.NewService(reminderRepo, emailSender, cfg.Email.ReportNotificationEmail)
	go reminderNotification.StartWorker(workerCtx, notificationService, time.Minute, "America/Mexico_City")
	log.Info("reminder notification worker started",
		zap.Bool("smtp_configured", emailSender.IsConfigured()),
		zap.String("smtp_host", cfg.Email.SMTPHost),
		zap.String("smtp_port", cfg.Email.SMTPPort),
		zap.Bool("smtp_user_configured", cfg.Email.SMTPUser != ""),
		zap.Bool("smtp_pass_configured", cfg.Email.SMTPPass != ""),
		zap.Bool("smtp_from_configured", cfg.Email.SMTPFrom != ""),
		zap.Bool("fallback_recipient_configured", cfg.Email.ReportNotificationEmail != ""),
		zap.String("timezone", "America/Mexico_City"),
		zap.Duration("interval", time.Minute),
	)
	if !emailSender.IsConfigured() {
		log.Warn("reminder notification emails will fail until smtp config is complete",
			zap.Bool("smtp_host_configured", cfg.Email.SMTPHost != ""),
			zap.Bool("smtp_port_configured", cfg.Email.SMTPPort != ""),
			zap.Bool("smtp_user_configured", cfg.Email.SMTPUser != ""),
			zap.Bool("smtp_pass_configured", cfg.Email.SMTPPass != ""),
			zap.Bool("smtp_from_configured", cfg.Email.SMTPFrom != ""),
		)
	}

	deps := api.Dependencies{
		DB:  db,
		Cfg: *cfg,
	}

	app := api.NewRouter(deps)

	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Info("server listening", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	stopReminderWorker()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited")
}
