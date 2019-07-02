package models

import (
	"fmt"
)

// Log levels
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warning"
	ErrorLevel = "error"
)

// Event log event.
type Event struct {
	AppName   string `json:"app,omitempty"`
	Version   string `json:"version,omitempty"`
	SessionID string `json:"sessionId,omitempty"`
	ClientID  string `json:"clientId,omitempty"`
	Message   string `json:"message,omitempty"`
	Level     string `json:"level,omitempty"`
}

func (e *Event) String() string {
	return fmt.Sprintf("Event(level=%s, app=%s, version=%s, sessionId=%s, clientId=%s, message=%s)",
		e.Level,
		e.AppName,
		e.Version,
		e.SessionID,
		e.ClientID,
		e.Message)
}
