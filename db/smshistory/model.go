package smshistory

import "time"

// SmsHistory represents an sms sent at a timestamp to a phonenumber
type SmsHistory struct {
	Phonenumber string
	CreatedAt   time.Time
}

// New creates a new SmsHistory
func New(phonenumber string) *SmsHistory {
	return &SmsHistory{
		Phonenumber: phonenumber,
		CreatedAt:   time.Now(),
	}
}
