package iyoid

import (
	"errors"
	"net/http"

	"github.com/itsyouonline/identityserver/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoIdentifierCollectionName = "iyoids"
	maxIdentifiers                = 25
)

var (
	// ErrIDLimitReached indicates that the max amount of ids has been reached
	ErrIDLimitReached = errors.New("Max amount of iyoids reached for this username and azp")
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:    []string{"username", "azp"},
		Unique: true,
	}

	db.EnsureIndex(mongoIdentifierCollectionName, index)

	index = mgo.Index{
		Key:    []string{"iyoids"},
		Unique: true,
	}

	db.EnsureIndex(mongoIdentifierCollectionName, index)

	index = mgo.Index{
		Key:    []string{"iyoids", "azp"},
		Unique: true,
	}

	db.EnsureIndex(mongoIdentifierCollectionName, index)
}

// Manager represents the database session
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

func (m *Manager) getIdentifierCollection() *mgo.Collection {
	return db.GetCollection(m.session, mongoIdentifierCollectionName)
}

// GetByIDAndAZP gets an Identifier object by IyoID.
func (m *Manager) GetByIDAndAZP(iyoid, azp string) (*Identifier, error) {
	var idObj Identifier

	if err := m.getIdentifierCollection().Find(bson.M{"iyoids": iyoid, "azp": azp}).One(&idObj); err != nil {
		return nil, err
	}

	return &idObj, nil
}

// GetByUsernameAndAZP returns the Identifier object for this username and azp combo
func (m *Manager) GetByUsernameAndAZP(username, azp string) (*Identifier, error) {
	var idObj Identifier

	if err := m.getIdentifierCollection().Find(bson.M{"username": username, "azp": azp}).One(&idObj); err != nil {
		return nil, err
	}

	return &idObj, nil
}

// UpsertIdentifier adds a new iyoid to a mapping or creates a new mapping
func (m *Manager) UpsertIdentifier(username, azp, iyoid string) error {
	// Count the amount of iyoids we already have first
	idObj, err := m.GetByUsernameAndAZP(username, azp)
	if err != nil && !db.IsNotFound(err) {
		return err
	}
	if db.IsNotFound(err) {
		idObj = &Identifier{}
	}
	if len(idObj.IyoIDs) >= maxIdentifiers {
		return ErrIDLimitReached
	}

	_, err = m.getIdentifierCollection().Upsert(bson.M{"username": username, "azp": azp}, bson.M{"$push": bson.M{"iyoids": iyoid}})
	return err
}
