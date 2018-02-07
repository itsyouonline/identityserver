package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type UserOrganizations struct {
	Member []string `json:"member" validate:"nonzero"`
	Owner  []string `json:"owner" validate:"nonzero"`
}

func (s UserOrganizations) Validate() error {

	return validator.Validate(s)
}
