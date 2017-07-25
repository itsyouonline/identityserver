package user

import (
	"errors"
	"regexp"

	"github.com/itsyouonline/identityserver/db"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

type EmailAddress struct {
	EmailAddress string `json:"emailaddress" validate:"max=100"`
	Label        string `json:"label" validate:"regexp=^[a-zA-Z\d\-_\s]{2,50}$"`
}

type PublicKey struct {
	PublicKey string `json:"publickey"`
	Label     string `json:"label" validate:"regexp=^[a-zA-Z\d\-_\s]{2,50}$"`
}

// Avatar represents an avatar for a user. It is identified by a label, and stirng
// contains a link to the source
type Avatar struct {
	Label  string `json:"label"`
	Source string `json:"source"`
}

type User struct {
	ID             bson.ObjectId         `json:"-" bson:"_id,omitempty"`
	Addresses      []Address             `json:"addresses"`
	BankAccounts   []BankAccount         `json:"bankaccounts"`
	EmailAddresses []EmailAddress        `json:"emailaddresses"`
	Expire         db.DateTime           `json:"-" bson:"expire,omitempty"`
	Facebook       FacebookAccount       `json:"facebook"`
	Github         GithubAccount         `json:"github"`
	Phonenumbers   []Phonenumber         `json:"phonenumbers"`
	DigitalWallet  []DigitalAssetAddress `json:"digitalwallet"`
	PublicKeys     []PublicKey           `json:"publicKeys"`
	Username       string                `json:"username" validate:"min=2,max=30,regexp=^[a-z0-9]{2,30}$"`
	Firstname      string                `json:"firstname"`
	Lastname       string                `json:"lastname"`
	Avatars        []Avatar              `json:"avatars"`
}

func (u *User) GetEmailAddressByLabel(label string) (email EmailAddress, err error) {
	for _, email = range u.EmailAddresses {
		if email.Label == label {
			return
		}
	}
	err = errors.New("Could not find EmailAddress with Label " + email.Label)
	return
}

func (u *User) GetPhonenumberByLabel(label string) (phonenumber Phonenumber, err error) {
	for _, phonenumber = range u.Phonenumbers {
		if phonenumber.Label == label {
			return
		}
	}
	err = errors.New("Could not find Phonenumber with Label " + phonenumber.Label)
	return
}

func (u *User) GetBankAccountByLabel(label string) (bankaccount BankAccount, err error) {
	for _, bankaccount = range u.BankAccounts {
		if bankaccount.Label == label {
			return
		}
	}
	err = errors.New("Could not find Phonenumber with Label " + bankaccount.Label)
	return
}

func (u *User) GetAddressByLabel(label string) (address Address, err error) {
	for _, address = range u.Addresses {
		if address.Label == label {
			return
		}
	}
	err = errors.New("Could not find Phonenumber with Label " + address.Label)
	return
}

func (u *User) GetDigitalAssetAddressByLabel(label string) (walletAddress DigitalAssetAddress, err error) {
	for _, walletAddress = range u.DigitalWallet {
		if walletAddress.Label == label {
			return
		}
	}
	err = errors.New("Could not find DigitalAssetAddress with Label " + walletAddress.Label)
	return
}

// GetPublicKeyByLabel Gets the public key associated with this label
func (u *User) GetPublicKeyByLabel(label string) (publicKey PublicKey, err error) {
	for _, publicKey = range u.PublicKeys {
		if publicKey.Label == label {
			return
		}
	}
	err = errors.New("Could not find PublicKey with label " + label)
	return
}

// GetAvatarByLabel gets the avatar associated with this label
func (u *User) GetAvatarByLabel(label string) (avatar Avatar, err error) {
	for _, avatar = range u.Avatars {
		if avatar.Label == label {
			return
		}
	}
	err = errors.New("Could not find Avatar with Label " + label)
	return
}

func ValidateUsername(username string) bool {
	regex, _ := regexp.Compile(`^[a-z\d\-_\s]{2,30}$`)
	matches := regex.FindAllString(username, 2)
	return len(matches) == 1
}

func ValidatePhoneNumber(phoneNumber string) bool {
	regex := regexp.MustCompile(`^\+[0-9 ]*$`)
	return regex.MatchString(phoneNumber)
}

func ValidateEmailAddress(emailAddress string) bool {
	regex := regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$")
	return regex.MatchString(emailAddress) && len(emailAddress) <= 100
}

func (p PublicKey) Validate() bool {
	return validator.Validate(p) == nil && regexp.MustCompile(`^[a-zA-Z\d\-_\s]{2,50}$`).MatchString(p.Label)
}

func (e EmailAddress) Validate() bool {

	return validator.Validate(e) == nil &&
		regexp.MustCompile(`^[a-zA-Z\d\-_\s]{2,50}$`).MatchString(e.Label) && ValidateEmailAddress(e.EmailAddress)
}
