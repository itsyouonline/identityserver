package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type TwoFAMethods struct {
	Sms  []Phonenumber `json:"sms" validate:"nonzero"`
	Totp bool          `json:"totp"`
}

func (s TwoFAMethods) Validate() error {

	return validator.Validate(s)
}
