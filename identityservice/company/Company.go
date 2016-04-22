package company

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Company struct {
	Id            bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Expire        time.Time     `json:"expire"`
	Globalid      string        `json:"globalid"`
	Info          []string      `json:"info"`
	Organizations []string      `json:"organizations"`
	PublicKeys    []string      `json:"publicKeys"`
	Taxnr         string        `json:"taxnr"`
}

// IsValid performs basic validation on the content of a company's fields
func (c *Company) IsValid() (valid bool) {
	valid = true
	globalIDLength := len(c.Globalid)
	valid = valid && (globalIDLength >= 3) && (globalIDLength <= 150)
	return
}
