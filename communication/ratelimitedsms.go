package communication

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/smshistory"
)

var (
	// ErrMaxSMS indicates that the sms sending to this number is being rate limited
	ErrMaxSMS = errors.New("Reached the max amount of sms for this phone number, try again later")
)

// RateLimitedSMSService wraps an SMS service and applies rate limiting based on the phone number
type RateLimitedSMSService struct {
	actualService SMSService
	window        time.Duration
	maxSMS        int
}

// NewRateLimitedSMSService rates limit an existing sms service to the defined rate
func NewRateLimitedSMSService(window int, maxSMS int, actualService SMSService) SMSService {
	return &RateLimitedSMSService{
		actualService: actualService,
		window:        time.Duration(int(time.Second) * window),
		maxSMS:        maxSMS,
	}
}

// Send checksif the message can be send according to the rate limiting rules, and then
// deligates the acutal sender to the wrapped service. This uses a sliding window approach
func (s *RateLimitedSMSService) Send(phonenumber string, message string) (err error) {
	conn := db.NewSession()
	if conn == nil {
		return errors.New("Failed to get DB connection")
	}

	mgr := smshistory.NewManager(conn)

	sendCount, err := mgr.CountSMSHistorySince(phonenumber, time.Now().Add(-s.window))
	if err != nil && !db.IsNotFound(err) {
		log.Error("Failed to count sms history: ", err)
		return err
	}

	// if we've send more or equal than max sms count msges already return an error
	if sendCount >= s.maxSMS {
		log.Info("Rate limiting sms sending to ", phonenumber, ", already sent ", sendCount, " SMS in the last ", s.window)
		return ErrMaxSMS
	}

	// Add log entry
	record := smshistory.New(phonenumber)
	if err = mgr.AddSMSHistory(record); err != nil {
		return err
	}

	// Send the actual sms
	return s.actualService.Send(phonenumber, message)
}
