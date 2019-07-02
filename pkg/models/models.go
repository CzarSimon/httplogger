package models

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
