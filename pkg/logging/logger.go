package logging

import (
	"context"
	"sync"

	"github.com/jklaw90/shinfo/pkg/config"
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
	if defaultLogger != nil {
		return defaultLogger
	}
	defaultLoggerOnce.Do(func() {
		config := zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.Level(zapcore.DebugLevel))
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		l, _ := config.Build()
		defaultLogger = l.Sugar()
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

func NewLogger(cfg config.Provider) *zap.SugaredLogger {
	var l *zap.Logger
	defaultLoggerOnce.Do(func() {
		if cfg.GetBool("logging.json") {
			config := zap.NewProductionConfig()
			config.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.GetInt("logging.level")))
			l, _ = config.Build()
		} else {
			config := zap.NewDevelopmentConfig()
			config.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.GetInt("logging.level")))
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			l, _ = config.Build()
		}
	})
	defaultLogger = l.Sugar()

	return defaultLogger
}
