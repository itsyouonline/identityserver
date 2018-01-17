package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type userview struct {
	Addresses               []Address             `json:"addresses" validate:"nonzero"`
	Avatar                  []Avatar              `json:"avatar" validate:"nonzero"`
	Bankaccounts            []BankAccount         `json:"bankaccounts" validate:"nonzero"`
	Digitalwallet           []DigitalAssetAddress `json:"digitalwallet" validate:"nonzero"`
	Emailaddresses          []EmailAddress        `json:"emailaddresses" validate:"nonzero"`
	Facebook                FacebookAccount       `json:"facebook,omitempty"`
	Firstname               string                `json:"firstname" validate:"nonzero"`
	Github                  GithubAccount         `json:"github,omitempty"`
	Lastname                string                `json:"lastname" validate:"nonzero"`
	Organizations           []string              `json:"organizations" validate:"nonzero"`
	Ownerof                 Ownerof               `json:"ownerof" validate:"nonzero"`
	Phonenumbers            []Phonenumber         `json:"phonenumbers" validate:"nonzero"`
	PublicKeys              []PublicKey           `json:"publicKeys,omitempty"`
	Username                string                `json:"username" validate:"nonzero"`
	Validatedemailaddresses []EmailAddress        `json:"validatedemailaddresses" validate:"nonzero"`
	Validatedphonenumbers   []Phonenumber         `json:"validatedphonenumbers" validate:"nonzero"`
}

func (s userview) Validate() error {

	return validator.Validate(s)
}
