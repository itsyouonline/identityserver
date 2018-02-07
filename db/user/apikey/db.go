package apikey

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoCollectionName = "apikeys"
)

//Manager is used to store users
type Manager struct {
	session *mgo.Session
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session: session,
	}
}

func (m *Manager) getCollection() *mgo.Collection {
	return db.GetCollection(m.session, mongoCollectionName)
}

//Save ApplicationAPIKey
func (m *Manager) Save(apikey *APIKey) (err error) {
	if apikey.ID == "" {
		// New Doc!
		apikey.ID = bson.NewObjectId()
		err = m.getCollection().Insert(apikey)
		return
	}
	_, err = m.getCollection().UpsertId(apikey.ID, apikey)
	return
}

func (m *Manager) GetByUsernameAndLabel(username string, label string) (apikey *APIKey, err error) {
	apikey = &APIKey{}
	err = m.getCollection().Find(bson.M{"username": username, "label": label}).One(apikey)
	if err == mgo.ErrNotFound {
		err = nil
	} else if err != nil {
		apikey = nil
	}
	return
}

// Exists checks if an api key for a user exists
func (m *Manager) Exists(username string, applicationid string) (bool, error) {
	count, err := m.getCollection().Find(bson.M{"username": username, "applicationid": applicationid}).Count()

	return count == 1, err
}

func (m *Manager) GetByApplicationAndSecret(applicationid string, secret string) (apikey *APIKey, err error) {
	apikey = &APIKey{}
	err = m.getCollection().Find(bson.M{"applicationid": applicationid, "apikey": secret}).One(apikey)
	return
}

func (m *Manager) GetByUser(username string) (apikeys []APIKey, err error) {
	err = m.getCollection().Find(bson.M{"username": username}).All(&apikeys)
	return
}

//Delete ApplicationAPIKey
func (m *Manager) Delete(username string, label string) (err error) {
	_, err = m.getCollection().RemoveAll(bson.M{"username": username, "label": label})
	return
}
