package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type TOTPSecret struct {
	Totpcode   string `json:"totpcode" validate:"nonzero"`
	Totpsecret string `json:"totpsecret" validate:"nonzero"`
}

func (s TOTPSecret) Validate() error {

	return validator.Validate(s)
}
