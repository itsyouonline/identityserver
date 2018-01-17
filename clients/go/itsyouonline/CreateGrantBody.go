package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type CreateGrantBody struct {
	Grant    Grant  `json:"grant" validate:"nonzero"`
	Username string `json:"username" validate:"nonzero"`
}

func (s CreateGrantBody) Validate() error {

	return validator.Validate(s)
}
