package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type UpdateGrantBody struct {
	Newgrant Grant  `json:"newgrant" validate:"nonzero"`
	Oldgrant Grant  `json:"oldgrant" validate:"nonzero"`
	Username string `json:"username" validate:"nonzero"`
}

func (s UpdateGrantBody) Validate() error {

	return validator.Validate(s)
}
