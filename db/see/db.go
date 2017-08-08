package see

import (
	"net/http"

	"gopkg.in/mgo.v2"
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
	qry := bson.M{"username": username}
	err = m.collection.Find(qry).Sort("-versions.creationdate").All(&seeObjects)
	if seeObjects == nil {
		seeObjects = []See{}
	}
	return
}

// GetSeeObjectsByOrganization returns all see object for a specific username / organization
func (m *Manager) GetSeeObjectsByOrganization(username string, globalID string) (seeObjects []See, err error) {
	qry := bson.M{"username": username, "globalid": globalID}
	err = m.collection.Find(qry).Sort("-versions.creationdate").All(&seeObjects)
	if seeObjects == nil {
		seeObjects = []See{}
	}
	return
}

// GetSeeObject returns a see object
func (m *Manager) GetSeeObject(username string, globalID string, uniqueID string) (seeObject *See, err error) {
	qry := bson.M{"username": username, "globalid": globalID, "uniqueid": uniqueID}
	err = m.collection.Find(qry).One(&seeObject)
	return
}

// Create an object
func (m *Manager) Create(see *See) error {
	see.ID = bson.NewObjectId()
	err := m.collection.Insert(see)
	if mgo.IsDup(err) {
		return db.ErrDuplicate
	}
	return err
}

// AddVersion adds a new version to the object
func (m *Manager) AddVersion(username string, globalID string, uniqueID string, seeVersion *SeeVersion) error {
	qry := bson.M{"username": username, "globalid": globalID, "uniqueid": uniqueID}
	return m.collection.Update(qry, bson.M{"$push": bson.M{"versions": seeVersion}})
}

// Update adds a signature to an existing version
func (m *Manager) Update(see *See) error {
	return m.collection.UpdateId(see.ID, see)
}
