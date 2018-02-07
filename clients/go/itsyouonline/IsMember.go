package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type IsMember struct {
	IsMember bool `json:"IsMember"`
}

func (s IsMember) Validate() error {

	return validator.Validate(s)
}
