package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type OrganizationsGlobalidLogoPutReqBody struct {
	Logo string `json:"logo" validate:"nonzero"`
}

func (s OrganizationsGlobalidLogoPutReqBody) Validate() error {

	return validator.Validate(s)
}
