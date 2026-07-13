package infra

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"gin-layout/config"
	"gin-layout/internal/domain"
)

type Logger = zerolog.Logger

func NewLogger(cfg *config.LogConfig) *Logger {
	output, err := openOutput(cfg.OutputPath)
	if err != nil {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	}

	zerolog.TimeFieldFormat = time.DateTime
	zerolog.SetGlobalLevel(parseLevel(cfg.Level))

	l := zerolog.New(output).With().Timestamp().Logger()
	if cfg.Format == "console" {
		l = l.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime})
	}
	return &l
}

func DefaultLogger() *Logger {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &l
}

func LogFromContext(ctx context.Context, base *Logger) zerolog.Logger {
	logger := zerolog.Ctx(ctx)
	if logger == nil {
		if base != nil {
			logger = base
		} else {
			logger = DefaultLogger()
		}
	}

	logCtx := logger.With()
	if requestID, ok := domain.RequestIDFromContext(ctx); ok {
		logCtx = logCtx.Str(domain.RequestIDKey, requestID)
	}
	if user, ok := domain.CurrentUserFromContext(ctx); ok {
		logCtx = logCtx.Int64(domain.UserIDKey, user.UserID).Str(domain.UserKey, user.Account)
	}

	return logCtx.Logger()
}

func openOutput(outputPath string) (io.Writer, error) {
	if outputPath == "" {
		return os.Stdout, nil
	}
	return os.OpenFile(outputPath+"/app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
