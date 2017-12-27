package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type PhoneNumberValidation struct {
	Validationkey string `json:"validationkey" validate:"nonzero"`
}

func (s PhoneNumberValidation) Validate() error {

	return validator.Validate(s)
}
