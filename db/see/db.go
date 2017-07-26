package see

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoCollectionName = "see"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:      []string{"username", "globalid", "uniqueid"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(mongoCollectionName, index)
}

//Manager is used to store users
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoCollectionName)
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: getCollection(session),
	}
}

// GetSeeObjects returns all see object for a specific username
func (m *Manager) GetSeeObjects(username string) (seeObjects []See, err error) {
	var see See
	see.Globalid = "test"
	see.Username = "ruben_1"
	see.Uniqueid = "test"
	see.Link = "https://github.com/itsyouonline/identityserver/issues/547"
	m.SaveSeeObject(see)

	qry := bson.M{"username": username}
	err = m.collection.Find(qry).All(&seeObjects)
	if seeObjects == nil {
		seeObjects = []See{}
	}
	return
}

// GetSeeObjectsByOrganization returns all see object for a specific username / organization
func (m *Manager) GetSeeObjectsByOrganization(username string, globalID string) (seeObjects []See, err error) {
	qry := bson.M{"username": username, "globalid": globalID}
	err = m.collection.Find(qry).All(&seeObjects)
	if seeObjects == nil {
		seeObjects = []See{}
	}
	return
}

// GetSeeObject returns a see object
func (m *Manager) GetSeeObject(username string, globalID string, uniqueID string) (seeObject See, err error) {
	qry := bson.M{"username": username, "globalid": globalID, "uniqueid": uniqueID}
	err = m.collection.Find(qry).One(&seeObject)
	return
}

// SaveSeeObject saves a new or updates an existing
func (m *Manager) SaveSeeObject(see See) error {

	if see.ID == "" {
		// New Doc!
		see.ID = bson.NewObjectId()
		err := m.collection.Insert(see)
		return err
	}

	_, err := m.collection.UpsertId(see.ID, see)

	return err
}
