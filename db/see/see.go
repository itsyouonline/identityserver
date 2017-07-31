package see

import (
	"time"

	"github.com/itsyouonline/identityserver/db"
	"gopkg.in/mgo.v2/bson"
)

type SeeVersion struct {
	Version                  int          `json:"version" bson:"-"`
	Category                 string       `json:"category"`
	Link                     string       `json:"link"`
	ContentType              string       `json:"content_type"`
	MarkdownShortDescription string       `json:"markdown_short_description"`
	MarkdownFullDescription  string       `json:"markdown_full_description"`
	CreationDate             *db.DateTime `json:"creation_date"`
	StartDate                *db.DateTime `json:"start_date,omitempty" bson:"startdate,omitempty"`
	EndDate                  *db.DateTime `json:"end_date,omitempty" bson:"enddate,omitempty"`
	KeyStoreLabel            string       `json:"keystore_label"`
	Signature                string       `json:"signature"`
}

type See struct {
	ID       bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Username string        `json:"username"`
	Globalid string        `json:"globalid"`
	Uniqueid string        `json:"uniqueid"`
	Versions []SeeVersion  `json:"versions"`
}

type SeeView struct {
	Username                 string       `json:"username"`
	Globalid                 string       `json:"globalid"`
	Uniqueid                 string       `json:"uniqueid"`
	Version                  int          `json:"version"`
	Category                 string       `json:"category"`
	Link                     string       `json:"link"`
	ContentType              string       `json:"content_type"`
	MarkdownShortDescription string       `json:"markdown_short_description"`
	MarkdownFullDescription  string       `json:"markdown_full_description"`
	CreationDate             *db.DateTime `json:"creation_date"`
	StartDate                *db.DateTime `json:"start_date,omitempty"`
	EndDate                  *db.DateTime `json:"end_date,omitempty"`
	KeyStoreLabel            string       `json:"keystore_label"`
	Signature                string       `json:"signature"`
}

func (s *SeeView) Validate() bool {
	return len(s.Uniqueid) < 100 && len(s.Uniqueid) > 0 &&
		len(s.Category) < 100 && len(s.Category) > 0 &&
		len(s.Link) > 0 &&
		len(s.MarkdownShortDescription) > 0 &&
		len(s.MarkdownFullDescription) > 0
}

func (s *SeeView) ConvertToSeeVersion() *SeeVersion {
	now := db.DateTime(time.Now())
	seeVersion := SeeVersion{}
	seeVersion.Category = s.Category
	seeVersion.Link = s.Link
	seeVersion.ContentType = s.ContentType
	seeVersion.MarkdownShortDescription = s.MarkdownShortDescription
	seeVersion.MarkdownFullDescription = s.MarkdownFullDescription
	seeVersion.CreationDate = &now
	seeVersion.StartDate = s.StartDate
	seeVersion.EndDate = s.EndDate
	seeVersion.KeyStoreLabel = s.KeyStoreLabel
	seeVersion.Signature = s.Signature
	return &seeVersion
}

func (s *See) ConvertToSeeView(version int) *SeeView {
	seeView := SeeView{}
	seeView.Username = s.Username
	seeView.Globalid = s.Globalid
	seeView.Uniqueid = s.Uniqueid
	seeView.Version = version
	seeView.Category = s.Versions[version-1].Category
	seeView.Link = s.Versions[version-1].Link
	seeView.ContentType = s.Versions[version-1].ContentType
	seeView.MarkdownShortDescription = s.Versions[version-1].MarkdownShortDescription
	seeView.MarkdownFullDescription = s.Versions[version-1].MarkdownFullDescription
	seeView.CreationDate = s.Versions[version-1].CreationDate
	seeView.StartDate = s.Versions[version-1].StartDate
	seeView.EndDate = s.Versions[version-1].EndDate
	seeView.KeyStoreLabel = s.Versions[version-1].KeyStoreLabel
	seeView.Signature = s.Versions[version-1].Signature
	return &seeView
}
