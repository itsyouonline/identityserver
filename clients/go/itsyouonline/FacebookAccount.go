package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type FacebookAccount struct {
	Id      string `json:"id" validate:"nonzero"`
	Link    string `json:"link" validate:"nonzero"`
	Name    string `json:"name" validate:"nonzero"`
	Picture string `json:"picture" validate:"nonzero"`
}

func (s FacebookAccount) Validate() error {

	return validator.Validate(s)
}
