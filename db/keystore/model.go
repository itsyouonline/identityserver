package keystore

import (
	"regexp"

	"github.com/itsyouonline/identityserver/db"
)

type KeyStoreKey struct {
	Key      string  `json:"key"`
	Globalid string  `json:"globalid"`
	Username string  `json:"username"`
	Label    string  `json:"label"`
	KeyData  KeyData `json:"keydata"`
}

type KeyData struct {
	TimeStamp db.DateTime `json:"timestamp"`
	Comment   string      `json:"comment"`
	Algorithm string      `json:"algorithm"`
}

func (key *KeyStoreKey) Validate() bool {
	return regexp.MustCompile(`^[a-zA-Z\d\-_\s]{2,50}$`).MatchString(key.Label) &&
		len(key.KeyData.Comment) < 100 && // comment is optional, max 100 bytes in char values
		len(key.KeyData.Algorithm) < 50 &&
		len(key.KeyData.Algorithm) > 0 &&
		len(key.Key) > 0
}
