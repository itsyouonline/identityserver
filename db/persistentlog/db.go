package persistentlog

import (
	"net/http"

	"github.com/itsyouonline/identityserver/db"
	mgo "gopkg.in/mgo.v2"
)

const (
	mongoPersistenLogCollectionName = "persistentlogs"
)

// InitModels intializes the mongo models
func InitModels() {
	index := mgo.Index{
		Key:    []string{"key"},
		Unique: false,
	}

	db.EnsureIndex(mongoPersistenLogCollectionName, index)
}

// Manager is used to store logs
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: db.GetCollection(session, mongoPersistenLogCollectionName),
	}
}

// SaveLog stores a new PersistentLog entry
func (m *Manager) SaveLog(log *PersistentLog) error {
	return m.collection.Insert(log)
}
