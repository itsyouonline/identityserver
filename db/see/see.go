package see

import (
	"github.com/itsyouonline/identityserver/db"
	"gopkg.in/mgo.v2/bson"
	validator "gopkg.in/validator.v2"
)

type See struct {
	ID                       bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Username                 string        `json:"username"`
	Globalid                 string        `json:"globalid"`
	Uniqueid                 string        `json:"uniqueid"`
	Link                     string        `json:"link"`
	Category                 string        `json:"category"`
	Version                  int           `json:"version"`
	ContentType              string        `json:"content_type"`
	MarkdownShortDescription string        `json:"markdown_short_description"`
	MarkdownFullDescription  string        `json:"markdown_full_description"`
	CreationDate             db.DateTime   `json:"creation_date"`
	StartDate                db.DateTime   `json:"start_date"`
	EndDate                  db.DateTime   `json:"end_date"`
	Signature                string        `json:"signature"`
}

func (s See) Validate() bool {
	return validator.Validate(s) == nil
}
