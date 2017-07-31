package user

import "github.com/itsyouonline/identityserver/db/user"

type Userview struct {
	Addresses               []user.Address             `json:"addresses"`
	BankAccounts            []user.BankAccount         `json:"bankaccounts"`
	DigitalWallet           []user.DigitalAssetAddress `json:"bankaccounts"`
	EmailAddresses          []user.EmailAddress        `json:"emailaddresses"`
	ValidatedEmailAddresses []user.EmailAddress        `json:"validatedemailaddresses"`
	Facebook                user.FacebookAccount       `json:"facebook"`
	Github                  user.GithubAccount         `json:"github"`
	Organizations           []string                   `json:"organizations"`
	Phonenumbers            []user.Phonenumber         `json:"phonenumbers"`
	ValidatedPhonenumbers   []user.Phonenumber         `json:"validatedphonenumbers"`
	PublicKeys              []user.PublicKey           `json:"publicKeys"`
	Username                string                     `json:"username"`
	Firstname               string                     `json:"firstname"`
	Lastname                string                     `json:"lastname"`
	OwnerOf                 user.OwnerOf               `json:"ownerof"`
	Avatars                 []user.Avatar              `json:"avatar"`
}
