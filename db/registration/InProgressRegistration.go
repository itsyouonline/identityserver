package registration

import (
	"time"
)

// InProgressRegistration represents a (partial) registration object
type InProgressRegistration struct {
	SessionKey         string
	CreatedAt          time.Time
	Firstname          string
	Lastname           string
	Phonenumber        string
	PhoneValidationKey string
	Email              string
	EmailValidationKey string
}

// New returns a new in progress registration object
func New(sessionKey string) *InProgressRegistration {
	return &InProgressRegistration{
		CreatedAt:  time.Now(),
		SessionKey: sessionKey,
	}
}
