package registration

import (
	"net/http"
	"time"

	"github.com/itsyouonline/identityserver/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoRegistrationsInProgressCollectionName = "registrationsinprogress"
)

// InitModels intializes the mongo models
func InitModels() {
	index := mgo.Index{
		Key:    []string{"sessionkey"},
		Unique: true,
	}

	db.EnsureIndex(mongoRegistrationsInProgressCollectionName, index)

	index = mgo.Index{
		Key:    []string{"emailvalidationkey"},
		Unique: false, // Uniqueness is enforced in the respective ongoing validation collection
	}

	db.EnsureIndex(mongoRegistrationsInProgressCollectionName, index)

	index = mgo.Index{
		Key:    []string{"phonevalidationkey"},
		Unique: false, // Uniqueness is enforced in the respective ongoing validation collection
	}

	db.EnsureIndex(mongoRegistrationsInProgressCollectionName, index)

	automaticExpiration := mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 60 * 60 * 24, // Remove after one day
		Background:  true,
	}

	db.EnsureIndex(mongoRegistrationsInProgressCollectionName, automaticExpiration)
}

// Manager is used to store organizations
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: db.GetCollection(session, mongoRegistrationsInProgressCollectionName),
	}
}

// UpsertRegisteringUser creates a new or updates an existing entry in the db for a user currenly registering
func (m *Manager) UpsertRegisteringUser(ipr *InProgressRegistration) error {
	selector := bson.M{"sessionkey": ipr.SessionKey}
	_, err := m.collection.Upsert(selector, ipr)
	return err
}

// GetRegisteringUserBySessionKey returns a user object of an in progress registration
func (m *Manager) GetRegisteringUserBySessionKey(sessionKey string) (*InProgressRegistration, error) {
	var ipr InProgressRegistration

	err := m.collection.Find(bson.M{"sessionkey": sessionKey}).One(&ipr)

	return &ipr, err
}

// DeleteRegisteringUser deletes a registering user
func (m *Manager) DeleteRegisteringUser(sessionKey string) error {
	return m.collection.Remove(bson.M{"sessionkey": sessionKey})
}

// GetRegisteringUserByEmailKey looks up a registering user by an email validation key
// func (m *Manager) GetRegisteringUserByEmailKey(emailKey string) (*InProgressRegistration, error) {
// 	var ipr *InProgressRegistration

// 	err := m.collection.Find(bson.M{"emailvalidationkey": emailKey}).One(ipr)

// 	return ipr, err
// }

// GetRegisteringUserByPhoneKey looks up a registering user by a phone validation key
// func (m *Manager) GetRegisteringUserByPhoneKey(phonekey string) (*InProgressRegistration, error) {
// 	var ipr *InProgressRegistration

// 	err := m.collection.Find(bson.M{"phonevalidationkey": phonekey}).One(ipr)

// 	return ipr, err
// }
