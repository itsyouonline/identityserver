package validation

import "time"

//ValidatedPhonenumber is a record of a phonenumber for a user and when it is validated
type ValidatedPhonenumber struct {
	Username    string
	Phonenumber string
	CreatedAt   time.Time
}

type PhonenumberValidationInformation struct {
	Key         string
	SMSCode     string
	Username    string
	Phonenumber string
	Confirmed   bool
	CreatedAt   time.Time
}

type ValidatedEmailAddress struct {
	Username     string
	EmailAddress string
	CreatedAt    time.Time
}

type EmailAddressValidationInformation struct {
	Key          string
	Secret       string
	Username     string
	EmailAddress string
	Confirmed    bool
	CreatedAt    time.Time
}
