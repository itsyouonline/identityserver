package globalconfig

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoCollectionName = "globalconfig"
)

type GlobalConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// InitModels initialize models in mongo, if required
func InitModels() {
	// TODO: Use model tags to ensure indices/constraints
	index := mgo.Index{
		Key:      []string{"key"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(mongoCollectionName, index)
}

// Manager is used to store settings
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getConfigCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoCollectionName)
}

// NewManager creates and initializes a new Manager
func NewManager() *Manager {
	session := db.GetSession()
	return &Manager{
		session:    session,
		collection: getConfigCollection(session),
	}
}

// GetByGet return a config key/value from key name
func (m *Manager) GetByKey(key string) (*GlobalConfig, error) {
	var config GlobalConfig

	err := m.collection.Find(bson.M{"key": key}).One(&config)

	return &config, err
}

func (m *Manager) Exists(key string) (bool, error) {
	count, err := m.collection.Find(bson.M{"key": key}).Count()

	return count >= 1, err
}

// Insert a config key
func (m *Manager) Insert(c *GlobalConfig) error {
	err := m.collection.Insert(c)

	return err
}

// Delete a config key
func (m *Manager) Delete(key string) error {
	config, err := m.GetByKey(key)

	if err != nil {
		return err
	}

	return m.collection.Remove(config)
}
