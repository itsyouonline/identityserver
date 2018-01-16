package grants

import (
	"errors"
	"net/http"

	"github.com/itsyouonline/identityserver/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	grantCollectionName = "grants"
	maxGrants           = 50
)

var (
	// ErrGrantLimitReached indicates that an organization can not add anymore grants for a user
	ErrGrantLimitReached = errors.New("Max amount of grants reached for this user by this organization")
)

// InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:    []string{"username", "globalid"},
		Unique: true,
	}

	db.EnsureIndex(grantCollectionName, index)

	index = mgo.Index{
		Key:    []string{"grants"},
		Unique: true,
	}

	db.EnsureIndex(grantCollectionName, index)

	index = mgo.Index{
		Key:    []string{"grants", "globalid"},
		Unique: true,
	}

	db.EnsureIndex(grantCollectionName, index)
}

// Manager is used to store grants
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, grantCollectionName)
}

// NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: getCollection(session),
	}
}

// GetGrantsForUser gets all the saved grants for a user and globalID
func (m *Manager) GetGrantsForUser(username string, globalID string) (*SavedGrants, error) {
	var sg SavedGrants

	err := m.collection.Find(bson.M{"username": username, "globalid": globalID}).One(&sg)

	return &sg, err
}

// GetByGrant returns all SavedGrants where the given grant is in the list of grants
func (m *Manager) GetByGrant(grant Grant, globalID string) ([]SavedGrants, error) {
	var usersWithGrant []SavedGrants

	if err := m.collection.Find(bson.M{"grants": grant, "globalid": globalID}).All(&usersWithGrant); err != nil {
		return nil, err
	}

	return usersWithGrant, nil
}

// UpserGrant adds a new grent for a user by an organization
func (m *Manager) UpserGrant(username, globalID string, grant Grant) error {
	// Count the amount of grants we already have first
	sg, err := m.GetGrantsForUser(username, globalID)
	if err != nil && !db.IsNotFound(err) {
		return err
	}
	if db.IsNotFound(err) {
		sg = &SavedGrants{}
	}
	if len(sg.Grants) >= maxGrants {
		return ErrGrantLimitReached
	}

	_, err = m.collection.Upsert(bson.M{"username": username, "globalid": globalID}, bson.M{"$addToSet": bson.M{"grants": grant}})
	return err
}

// UpdateGrant updates an old grant to a new one
func (m *Manager) UpdateGrant(username, globalID string, oldgrant, newgrant Grant) error {
	// First remove the old grant
	err := m.collection.Update(bson.M{"username": username, "globalid": globalID}, bson.M{"$pull": bson.M{"grants": oldgrant}})
	if err != nil {
		return err
	}
	// Now insert the new one
	return m.collection.Update(bson.M{"username": username, "globalid": globalID}, bson.M{"$addToSet": bson.M{"grants": newgrant}})
}

// DeleteUserGrant removes a single grant from a user for an organization
func (m *Manager) DeleteUserGrant(username, globalID string, grant Grant) error {
	return m.collection.Update(bson.M{"username": username, "globalid": globalID}, bson.M{"$pull": bson.M{"grants": grant}})
}

// DeleteUserGrants removes all grants given to a user by an organization
func (m *Manager) DeleteUserGrants(username, globalID string) error {
	return m.collection.Remove(bson.M{"username": username, "globalid": globalID})
}

// DeleteOrgGrants remooves all grants given by an organization
func (m *Manager) DeleteOrgGrants(globalID string) error {
	_, err := m.collection.RemoveAll(bson.M{"globalid": globalID})
	return err
}
