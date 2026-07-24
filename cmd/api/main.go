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
	"github.com/agnathor/finances-go/internal/config"
	"github.com/agnathor/finances-go/internal/infrastructure/cache"
	"github.com/agnathor/finances-go/internal/infrastructure/database"
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

	rdb, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Warn("redis not available, continuing without cache", zap.Error(err))
	}
	if rdb != nil {
		defer rdb.Close()
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited")
}
