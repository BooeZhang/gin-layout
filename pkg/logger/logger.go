package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLog(level string, format string) {
	zerolog.TimeFieldFormat = time.DateTime
	lv, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid log level")
	}
	zerolog.SetGlobalLevel(lv)
	if format == "console" {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}
}

func SubLog(ctx *gin.Context) zerolog.Logger {
	requestID := ctx.GetString(XRequestIDKey)
	user := ctx.GetString(UserKey)
	subLog := log.With().Timestamp().Str("requestID", requestID).Str("user", user).Logger()
	return subLog
}
