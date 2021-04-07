package logging

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// https://github.com/google/exposure-notifications-server/blob/main/pkg/logging/logger.go
type contextKey string

const loggerKey = contextKey("logger")

var (
	defaultLogger     *zap.SugaredLogger
	defaultLoggerOnce sync.Once
)

func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger(Config{})
	})
	return defaultLogger
}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}

func NewLogger(cfg Config) *zap.SugaredLogger {
	var l *zap.Logger
	defaultLoggerOnce.Do(func() {
		if cfg.IsJson {
			config := zap.NewProductionConfig()
			config.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Level))
			l, _ = config.Build()
		} else {
			config := zap.NewDevelopmentConfig()
			config.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Level))
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			l, _ = config.Build()
		}
	})
	defaultLogger = l.Sugar()

	return defaultLogger
}
