package keystore

import (
	"net/http"

	"github.com/itsyouonline/identityserver/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoKeyStoreCollectionName = "keystore"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:      []string{"label", "username", "globalid"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(mongoKeyStoreCollectionName, index)
}

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

func (m *Manager) getKeyStoreCollection() *mgo.Collection {
	return db.GetCollection(m.session, mongoKeyStoreCollectionName)
}

// Create a new KeyStore key entry.
func (m *Manager) Create(key *KeyStoreKey) error {
	err := m.getKeyStoreCollection().Insert(key)
	if mgo.IsDup(err) {
		return db.ErrDuplicate
	}
	return err
}

func (m *Manager) ListKeyStoreKeys(username string, globalid string) ([]KeyStoreKey, error) {
	keys := make([]KeyStoreKey, 0)
	condition := []interface{}{
		bson.M{"globalid": globalid},
		bson.M{"username": username},
	}
	err := m.getKeyStoreCollection().Find(bson.M{"$and": condition}).All(&keys)
	return keys, err
}

func (m *Manager) GetKeyStoreKey(username string, globalid string, label string) (*KeyStoreKey, error) {
	qry := bson.M{
		"globalid": globalid,
		"username": username,
		"label":    label,
	}

	key := &KeyStoreKey{}
	err := m.getKeyStoreCollection().Find(qry).One(key)
	return key, err
}
