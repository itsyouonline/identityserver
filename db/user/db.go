package user

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"time"

	"github.com/itsyouonline/identityserver/db"
	"strings"
)

const (
	mongoUsersCollectionName          				= "users"
	mongoAvatarFileCollectionName     				= "avatarfiles"
	mongoAuthorizationsCollectionName 				= "authorizations"
	mongoJoinOrganizationInvitations  			    = "join-organization-invitations"
	mongoTotpCollectionName			  				= "totp"
	mongoValidatedAddressesCollectionName			= "validatedaddresses"
	mongoValidatedEmailAddressesCollectionName		= "validatedemailaddresses"
	mongoValidatedPhoneNumbersCollectionName		= "validatedphonenumbers"
	mongoOrganizationsCollectionName				= "organizations"
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

	avatarIndex := mgo.Index{
		Key:      []string{"hash"},
		Unique:   true,
		DropDups: true,
	}
	db.EnsureIndex(mongoAvatarFileCollectionName, avatarIndex)

	emailIndex := mgo.Index{
		Key: []string{"emailaddresses.emailaddress"},
	}
	db.EnsureIndex(mongoUsersCollectionName, emailIndex)
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

func (m *Manager) getAvatarFileCollection() *mgo.Collection {
	return db.GetCollection(m.session, mongoAvatarFileCollectionName)
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
	if user.Avatars == nil {
		user.Avatars = []Avatar{}
	}
	// Get the facebook and github profile pictures if they are linked to the account
	if user.Github.Avatar_url != "" {
		user.Avatars = append(user.Avatars, Avatar{Label: "github", Source: user.Github.Avatar_url})
	}
	if user.Facebook.Link != "" {
		user.Avatars = append(user.Avatars, Avatar{Label: "facebook", Source: user.Facebook.Picture})
	}
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
	if user.DigitalWallet == nil {
		user.DigitalWallet = []DigitalAssetAddress{}
	}
	if user.PublicKeys == nil {
		user.PublicKeys = []PublicKey{}
	}

	return &user, err
}

func (m *Manager) GetByEmailAddress(email string) (users []string, err error) {
	err = m.getUserCollection().Find(bson.M{"emailaddresses.emailaddress": email}).Distinct("username", &users)
	return
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

	// Check if the user is the only owner of any organization
	orgCollection := db.GetCollection(m.session, mongoOrganizationsCollectionName)
	var organizations []struct{}
	orgCollection.Find(bson.M{"owners": []string{u.Username}}).All(&organizations)
	if len(organizations) > 0 {
		return errors.New("only_owner")
	}

	// Remove the user from all organizations he is co-owner or member of
	orgCollection.UpdateAll(bson.M{}, bson.M{
		"$pull": bson.M{ "owners": u.Username, "members": u.Username },
	})

	// Remove all authorizations the user granted before for any organization
	_, err := m.getAuthorizationCollection().RemoveAll(bson.M{"username": u.Username})
	if err != nil {
		return err
	}

	// Delete all avatar files for the user
	for _, avatar := range u.Avatars {
		srcParts := strings.Split(avatar.Source, "/")
		err = m.RemoveAvatarFile(srcParts[len(srcParts)-1])
		if err != nil {
			return err
		}
	}

	// Delete all user's join organization invitations
	_, err = db.GetCollection(m.session, mongoJoinOrganizationInvitations).RemoveAll(bson.M{"username": u.Username})
	if err != nil {
		return err
	}

	// Delete all totp for the user
	_, err = db.GetCollection(m.session, mongoTotpCollectionName).RemoveAll(bson.M{"username": u.Username})
	if err != nil {
		return err
	}

	// Delete all validated addresses
	_, err = db.GetCollection(m.session, mongoValidatedAddressesCollectionName).RemoveAll(bson.M{"username": u.Username})
	if err != nil {
		return err
	}

	// Delete all validated phone numbers
	_, err = db.GetCollection(m.session, mongoValidatedPhoneNumbersCollectionName).RemoveAll(bson.M{"username": u.Username})
	if err != nil {
		return err
	}

	// Delete all validated email addresses
	_, err = db.GetCollection(m.session, mongoValidatedEmailAddressesCollectionName).RemoveAll(bson.M{"username": u.Username})
	if err != nil {
		return err
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

// SavePublicKey save or update public key along with its label
func (m *Manager) SavePublicKey(username string, key PublicKey) error {
	if err := m.RemovePublicKey(username, key.Label); err != nil {
		return err
	}
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$push": bson.M{"publickeys": key}})
}

// RemovePublicKey remove public key associated with label
func (m *Manager) RemovePublicKey(username string, label string) error {
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$pull": bson.M{"publickeys": bson.M{"label": label}}})
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

// GetOrganizationAuthorizations returns all authorizations for a specific organization
func (m *Manager) GetOrganizationAuthorizations(globalId string) (authorizations []Authorization, err error) {
	qry := bson.M{"grantedto": globalId}
	err = m.getAuthorizationCollection().Find(qry).All(&authorizations)
	if authorizations == nil {
		authorizations = []Authorization{}
	}
	return
}

//GetAuthorization returns the authorization for a specific organization, nil if no such authorization exists
func (m *Manager) GetAuthorization(username, organization string) (authorization *Authorization, err error) {
	err = m.getAuthorizationCollection().Find(bson.M{"username": username, "grantedto": organization}).One(&authorization)
	if err == mgo.ErrNotFound {
		err = nil
	} else if err != nil {
		authorization = nil
	}
	return
}

// FilterUsersWithAuthorizations returns the authorizations granted to an organization
func (m *Manager) FilterUsersWithAuthorizations(usernames []string, organization string) ([]string, error) {
	var auths []Authorization
	err := m.getAuthorizationCollection().Find(bson.M{"username": bson.M{"$in": usernames}, "grantedto": organization}).Select(bson.M{"username": 1}).All(&auths)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, auth := range auths {
		result = append(result, auth.Username)
	}
	return result, nil
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

// SaveAvatar saves a new or updates an existing avatar
func (m *Manager) SaveAvatar(username string, avatar Avatar) error {
	if err := m.RemoveAvatar(username, avatar.Label); err != nil {
		return err
	}
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$push": bson.M{"avatars": avatar}})
}

// RemoveAvatar removes an avatar
func (m *Manager) RemoveAvatar(username, label string) error {
	return m.getUserCollection().Update(
		bson.M{"username": username},
		bson.M{"$pull": bson.M{"avatars": bson.M{"label": label}}})
}

// AvatarFileExists checks if an avatarfile with a given hash already exists
func (m *Manager) AvatarFileExists(hash string) (bool, error) {
	count, err := m.getAvatarFileCollection().Find(bson.M{"hash": hash}).Count()

	return count >= 1, err
}

// GetAvatarFile gets the avatar file associated with a hash
func (m *Manager) GetAvatarFile(hash string) ([]byte, error) {
	var file struct {
		Hash string
		File []byte
	}

	err := m.getAvatarFileCollection().Find(bson.M{"hash": hash}).One(&file)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	return file.File, err
}

// SaveAvatarFile saves a new avatar file
func (m *Manager) SaveAvatarFile(hash string, file []byte) error {
	_, err := m.getAvatarFileCollection().Upsert(
		bson.M{"hash": hash},
		bson.M{"$set": bson.M{"file": file}})
	return err
}

// RemoveAvatarFile removes an avatar file
func (m *Manager) RemoveAvatarFile(hash string) error {
	return m.getAvatarFileCollection().Remove(bson.M{"hash": hash})
}
