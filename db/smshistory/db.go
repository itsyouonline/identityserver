package smshistory

import (
	"time"

	"github.com/itsyouonline/identityserver/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	smshistoryCollectionName = "smshistory"
)

// InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:         []string{"createdat"},
		Unique:      true,
		ExpireAfter: time.Duration(30*24*60*60) * time.Second, // one month
	}

	db.EnsureIndex(smshistoryCollectionName, index)
}

// Manager is used to store grants
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, smshistoryCollectionName)
}

// NewManager creates and initializes a new Manager
func NewManager(sess *mgo.Session) *Manager {
	return &Manager{
		session:    sess,
		collection: getCollection(sess),
	}
}

// AddSMSHistory adds SmsHistory to the database
func (m *Manager) AddSMSHistory(sh *SmsHistory) error {
	return m.collection.Insert(sh)
}

// CountSMSHistorySince counts the amount of sms sent to a phone number since a specific time
func (m *Manager) CountSMSHistorySince(phonenumber string, since time.Time) (int, error) {
	return m.collection.Find(bson.M{"createdat": bson.M{"$gte": since}}).Count()
}
