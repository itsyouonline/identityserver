package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type OrganizationLogo struct {
	Logo string `json:"logo" validate:"nonzero"`
}

func (s OrganizationLogo) Validate() error {

	return validator.Validate(s)
}
