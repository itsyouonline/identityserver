package persistentlog

import "time"

// PersistentLog represents a persistent log entry stored in the database
type PersistentLog struct {
	Timestamp time.Time
	Flow      FlowType
	Key       string
	Message   string
}

// FlowType represents the flow we are in when creating the log
type FlowType string

const (
	// RegistrationFlow is the FlowType for registration
	RegistrationFlow FlowType = "registration"
)

// New creates a new PersistentLog
func New(key string, flow FlowType, message string) *PersistentLog {
	return &PersistentLog{
		Timestamp: time.Now(),
		Flow:      flow,
		Key:       key,
		Message:   message,
	}
}
