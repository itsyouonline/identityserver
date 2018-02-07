package itsyouonline

import (
	"gopkg.in/validator.v2"
)

// Mapping between an Iyo ID, username and azp
type IyoID struct {
	Azp      string   `json:"azp" validate:"nonzero"`
	Iyoids   []string `json:"iyoids" validate:"nonzero"`
	Username string   `json:"username" validate:"nonzero"`
}

func (s IyoID) Validate() error {

	return validator.Validate(s)
}
