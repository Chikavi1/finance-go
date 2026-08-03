package remindernotification

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/agnathor/finances-go/pkg/logger"
)

func StartWorker(ctx context.Context, service Service, interval time.Duration, tzName string) {
	log := logger.Get()
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		log.Warn("failed to load reminder timezone, using local time", zap.String("timezone", tzName), zap.Error(err))
		loc = time.Local
	}

	run := func() {
		now := time.Now().In(loc)
		log.Info("starting reminder notification worker tick",
			zap.String("timestamp", now.Format(time.RFC3339)),
			zap.String("timezone", loc.String()),
		)
		sent, err := service.SendDue(ctx, now)
		if err != nil {
			log.Warn("reminder notification worker tick failed",
				zap.String("timestamp", now.Format(time.RFC3339)),
				zap.String("timezone", loc.String()),
				zap.Error(err),
			)
			return
		}
		log.Info("finished reminder notification worker tick",
			zap.String("timestamp", now.Format(time.RFC3339)),
			zap.String("timezone", loc.String()),
			zap.Int("sent_count", sent),
		)
	}

	run()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
