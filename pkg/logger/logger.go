package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(level string) {
	atomicLevel := zap.NewAtomicLevel()
	switch level {
	case "debug":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	default:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	if os.Getenv("APP_ENV") == "production" {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	writer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, writer, atomicLevel)

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func Get() *zap.Logger {
	if log == nil {
		Init("info")
	}
	return log
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
