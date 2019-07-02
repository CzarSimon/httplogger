package log

import (
	stdLog "log"

	"github.com/CzarSimon/httplogger/pkg/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	l, err := zap.NewProduction()
	if err != nil {
		stdLog.Fatalln("Failed to get zap.Logger", err)
	}

	logger = l
}

type logFn func(msg string, fields ...zapcore.Field)

// Log logs an event and records metrics for it.
func Log(e *models.Event) {
	log := selectLog(e)
	log(e.Message,
		zap.String("app", e.AppName),
		zap.String("version", e.Version),
		zap.String("sessionId", e.SessionID),
		zap.String("clientId", e.ClientID))
}

func selectLog(e *models.Event) logFn {
	switch e.Level {
	case models.DebugLevel:
		return logger.Debug
	case models.InfoLevel:
		return logger.Debug
	case models.WarnLevel:
		return logger.Warn
	case models.ErrorLevel:
		return logger.Error
	default:
		return logger.Info
	}
}
