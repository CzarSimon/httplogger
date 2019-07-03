package log_test

import (
	"testing"

	"github.com/CzarSimon/httplogger/pkg/log"
	"github.com/CzarSimon/httplogger/pkg/models"
)

func TestLog(t *testing.T) {
	event := &models.Event{
		Level:     "debug",
		AppName:   "app-1",
		Version:   "v1",
		SessionID: "session-id-1",
		ClientID:  "client-id-1",
		Message:   "debug test message",
	}
	log.Log(event)

	event = &models.Event{
		Level:     "info",
		AppName:   "app-1",
		Version:   "v1",
		SessionID: "session-id-1",
		ClientID:  "client-id-1",
		Message:   "info test message",
	}
	log.Log(event)

	event = &models.Event{
		Level:     "warning",
		AppName:   "app-1",
		Version:   "v1",
		SessionID: "session-id-1",
		ClientID:  "client-id-1",
		Message:   "warning test message",
	}
	log.Log(event)

	event = &models.Event{
		Level:      "error",
		AppName:    "app-1",
		Version:    "v1",
		SessionID:  "session-id-1",
		ClientID:   "client-id-1",
		Message:    "error test message",
		Stacktrace: "details on the error",
	}
	log.Log(event)

	event = &models.Event{
		Level:     "custom",
		AppName:   "app-1",
		Version:   "v1",
		SessionID: "session-id-1",
		ClientID:  "client-id-1",
		Message:   "custom test message",
	}
	log.Log(event)
}
