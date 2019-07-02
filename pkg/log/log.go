package log

import (
	stdLog "log"

	"github.com/CzarSimon/httplogger/pkg/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger       = setupLogger()
	eventsLogged = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_logged_total",
			Help: "The total number of logged events",
		},
		[]string{"app", "version", "level"},
	)
)

type logFn func(msg string, fields ...zapcore.Field)

// Log logs an event and records metrics for it.
func Log(e *models.Event) {
	log := selectLog(e)
	log(e.Message,
		zap.String("app", e.AppName),
		zap.String("version", e.Version),
		zap.String("sessionId", e.SessionID),
		zap.String("clientId", e.ClientID))

	eventsLogged.WithLabelValues(e.AppName, e.Version, e.Level).Inc()
}

func selectLog(e *models.Event) logFn {
	switch e.Level {
	case models.DebugLevel:
		return logger.Debug
	case models.InfoLevel:
		return logger.Info
	case models.WarnLevel:
		return logger.Warn
	case models.ErrorLevel:
		return logger.Error
	default:
		return logger.Info
	}
}

func setupLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true

	logger, err := cfg.Build()
	if err != nil {
		stdLog.Fatalln("Failed to get zap.Logger", err)
	}

	return logger
}
