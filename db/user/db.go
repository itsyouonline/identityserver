package user

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"time"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoUsersCollectionName          = "users"
	mongoAuthorizationsCollectionName = "authorizations"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:      []string{"username"},
		Unique:   true,
		DropDups: true,
	}

	db.EnsureIndex(mongoUsersCollectionName, index)

	// Removes users without valid 2 factor authentication after 3 days
	automaticUserExpiration := mgo.Index{
		Key:         []string{"expire"},
		ExpireAfter: time.Second * 3600 * 24 * 3,
		Background:  true,
	}
	db.EnsureIndex(mongoUsersCollectionName, automaticUserExpiration)
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

func (m *Manager) getUserCollection() *mgo.Collection {
	return db.GetCollection(m.session, mongoUsersCollectionName)
}

func (m *Manager) getAuthorizationCollection() *mgo.Collection {
	return db.GetCollection(m.session, mongoAuthorizationsCollectionName)
}

// Get user by ID.
func (m *Manager) Get(id string) (*User, error) {
	var user User

	objectID := bson.ObjectIdHex(id)

	if err := m.getUserCollection().FindId(objectID).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

//GetByName gets a user by it's username.
func (m *Manager) GetByName(username string) (*User, error) {
	var user User

	err := m.getUserCollection().Find(bson.M{"username": username}).One(&user)

	if user.Addresses == nil {
		user.Addresses = []Address{}
	}
	if user.BankAccounts == nil {
		user.BankAccounts = []BankAccount{}
	}
	if user.Phonenumbers == nil {
		user.Phonenumbers = []Phonenumber{}
	}
	if user.EmailAddresses == nil {
		user.EmailAddresses = []EmailAddress{}
	}

	return &user, err
}

//Exists checks if a user with this username already exists.
func (m *Manager) Exists(username string) (bool, error) {
	count, err := m.getUserCollection().Find(bson.M{"username": username}).Count()

	return count >= 1, err
}

// Save a user.
func (m *Manager) Save(u *User) error {
	// TODO: Validation!

	if u.ID == "" {
		// New Doc!
		u.ID = bson.NewObjectId()
		err := m.getUserCollection().Insert(u)
		return err
	}

	_, err := m.getUserCollection().UpsertId(u.ID, u)

	return err
}

// Delete a user.
func (m *Manager) Delete(u *User) error {
	if u.ID == "" {
		return errors.New("User not stored")
	}

	return m.getUserCollection().RemoveId(u.ID)
}

// SaveEmail save or update email along with its label
func (m *Manager) SaveEmail(username string, email EmailAddress) error {
	if err := m.RemoveEmail(username, email.Label); err != nil {
		return err
	}
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$push": bson.M{"emailaddresses": email}})
}

// RemoveEmail remove email associated with label
func (m *Manager) RemoveEmail(username string, label string) error {
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$pull": bson.M{"emailaddresses": bson.M{"label": label}}})
}

// SavePhone save or update phone along with its label
func (m *Manager) SavePhone(username string, phonenumber Phonenumber) error {
	if err := m.RemovePhone(username, phonenumber.Label); err != nil {
		return err
	}
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$push": bson.M{"phonenumbers": phonenumber}})
}

// RemovePhone remove phone associated with label
func (m *Manager) RemovePhone(username string, label string) error {
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$pull": bson.M{"phonenumbers": bson.M{"label": label}}})
}

// SaveVirtualCurrency save or update virtualcurrency along with its label
func (m *Manager) SaveVirtualCurrency(username string, currency DigitalAssetAddress) error {
	if err := m.RemoveVirtualCurrency(username, currency.Label); err != nil {
		return err
	}
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$push": bson.M{"digitalwallet": currency}})
}

// RemoveVirtualCurrency remove phone associated with label
func (m *Manager) RemoveVirtualCurrency(username string, label string) error {
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$pull": bson.M{"digitalwallet": bson.M{"label": label}}})
}

// SaveAddress save or update address
func (m *Manager) SaveAddress(username string, address Address) error {
	if err := m.RemoveAddress(username, address.Label); err != nil {
		return err
	}
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$push": bson.M{"addresses": address}})
}

// RemoveAddress remove address associated with label
func (m *Manager) RemoveAddress(username, label string) error {
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$pull": bson.M{"addresses": bson.M{"label": label}}})
}

// SaveBank save or update bank account
func (m *Manager) SaveBank(u *User, bank BankAccount) error {
	if err := m.RemoveBank(u, bank.Label); err != nil {
		return err
	}
	return m.getUserCollection().Update(
		bson.M{"username": u.Username},
		bson.M{"$push": bson.M{"bankaccounts": bank}})
}

// RemoveBank remove bank associated with label
func (m *Manager) RemoveBank(u *User, label string) error {
	return m.getUserCollection().Update(
		bson.M{"username": u.Username},
		bson.M{"$pull": bson.M{"bankaccounts": bson.M{"label": label}}})
}

func (m *Manager) UpdateGithubAccount(username string, githubaccount GithubAccount) (err error) {
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": bson.M{"github": githubaccount}})
	return
}

func (m *Manager) DeleteGithubAccount(username string) (err error) {
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": bson.M{"github": bson.M{}}})
	return
}

func (m *Manager) UpdateFacebookAccount(username string, facebookaccount FacebookAccount) (err error) {
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": bson.M{"facebook": facebookaccount}})
	return
}

func (m *Manager) DeleteFacebookAccount(username string) (err error) {
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": bson.M{"facebook": bson.M{}}})
	return
}

// GetAuthorizationsByUser returns all authorizations for a specific user
func (m *Manager) GetAuthorizationsByUser(username string) (authorizations []Authorization, err error) {
	err = m.getAuthorizationCollection().Find(bson.M{"username": username}).All(&authorizations)
	if authorizations == nil {
		authorizations = []Authorization{}
	}
	return
}

//GetAuthorization returns the authorization for a specific organization, nil if no such auhorization exists
func (m *Manager) GetAuthorization(username, organization string) (authorization *Authorization, err error) {
	authorization = &Authorization{}
	err = m.getAuthorizationCollection().Find(bson.M{"username": username, "grantedto": organization}).One(authorization)
	if err == mgo.ErrNotFound {
		err = nil
	} else if err != nil {
		authorization = nil
	}
	return
}

//UpdateAuthorization inserts or updates an authorization
func (m *Manager) UpdateAuthorization(authorization *Authorization) (err error) {
	_, err = m.getAuthorizationCollection().Upsert(bson.M{"username": authorization.Username, "grantedto": authorization.GrantedTo}, authorization)
	return
}

//DeleteAuthorization removes an authorization
func (m *Manager) DeleteAuthorization(username, organization string) (err error) {
	_, err = m.getAuthorizationCollection().RemoveAll(bson.M{"username": username, "grantedto": organization})
	return
}

//DeleteAllAuthorizations removes all authorizations from an organization
func (m *Manager) DeleteAllAuthorizations(organization string) (err error) {
	_, err = m.getAuthorizationCollection().RemoveAll(bson.M{"grantedto": organization})
	return err
}

func (u *User) getID() string {
	return u.ID.Hex()
}

func (m *Manager) UpdateName(username string, firstname string, lastname string) (err error) {
	values := bson.M{
		"firstname": firstname,
		"lastname":  lastname,
	}
	_, err = m.getUserCollection().UpdateAll(bson.M{"username": username}, bson.M{"$set": values})
	return
}

func (m *Manager) RemoveExpireDate(username string) (err error) {
	qry := bson.M{"username": username}
	values := bson.M{"expire": bson.M{}}
	_, err = m.getUserCollection().UpdateAll(qry, bson.M{"$set": values})
	return
}

func (m *Manager) GetPendingRegistrationsCount() (int, error) {
	qry := bson.M{
		"expire": bson.M{
			"$nin":    []interface{}{"", bson.M{}},
			"$exists": 1,
		},
	}
	return m.getUserCollection().Find(qry).Count()
}
