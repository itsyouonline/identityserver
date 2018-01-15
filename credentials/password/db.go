package password

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/credentials/password/keyderivation"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/tools"
)

var (
	// ErrInvalidPassword is returned when an invalid password is passed
	ErrInvalidPassword = errors.New("Invalid password")
)

const (
	mongoCollectionName           = "password"
	mongoCollectionNameResetToken = "passwordresetoken"
	passwordMinLength             = 6
)

type userPass struct {
	Username string
	Password string
}

type ResetToken struct {
	Username  string
	Token     string
	CreatedAt time.Time
}

//InitModels initializes models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:      []string{"username"},
		Unique:   true,
		DropDups: true,
	}
	db.EnsureIndex(mongoCollectionName, index)

	index = mgo.Index{
		Key:      []string{"username", "token"},
		Unique:   true,
		DropDups: true,
	}
	db.EnsureIndex(mongoCollectionNameResetToken, index)

	automaticExpiration := mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Minute * 10,
		Background:  true,
	}
	db.EnsureIndex(mongoCollectionNameResetToken, automaticExpiration)
}

//Manager stores and validates passwords
type Manager struct {
	session         *mgo.Session
	collection      *mgo.Collection
	tokencollection *mgo.Collection
}

func getPasswordCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoCollectionName)
}

func getPasswordResetTokenCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoCollectionNameResetToken)
}

//NewManager creates a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:         session,
		collection:      getPasswordCollection(session),
		tokencollection: getPasswordResetTokenCollection(session),
	}
}

//Validate checks the password for a specific username
func (pwm *Manager) Validate(username, password string) (bool, error) {
	var storedPassword userPass
	if err := pwm.collection.Find(bson.M{"username": username}).One(&storedPassword); err != nil {
		if err == mgo.ErrNotFound {
			log.Debug("No password found for this user")
			return false, nil
		}
		log.Debug(err)
		return false, err
	}
	match := keyderivation.Check(password, storedPassword.Password)
	return match, nil
}

// Save stores a password for a specific username.
func (pwm *Manager) Save(username, password string) error {
	//TODO: username and password validation
	passwordHash, err := keyderivation.Hash(password)
	if err != nil {
		log.Error("ERROR hashing password")
		log.Debug("ERROR hashing password: ", err)
		return errors.New("internal_error")
	}
	if len(password) < passwordMinLength {
		return errors.New("invalid_password")
	}
	storedPassword := userPass{Username: username, Password: passwordHash}

	_, err = pwm.collection.Upsert(bson.M{"username": username}, storedPassword)

	return err
}

// Check to see if a password is valid for a user.
func Check(password string) error {
	if len(password) < passwordMinLength {
		return ErrInvalidPassword
	}
	_, err := keyderivation.Hash(password)
	if err != nil {
		log.Debug("Failed to hash password: ", err)
		return err
	}
	return nil
}

// NewResetToken get new reset token
func (pwm *Manager) NewResetToken(username string) (token *ResetToken, err error) {
	tokenstring, err := tools.GenerateRandomString()
	if err != nil {
		return
	}
	token = &ResetToken{Token: tokenstring, CreatedAt: time.Now(), Username: username}
	return
}

// SaveResetToken save reset token
func (pwm *Manager) SaveResetToken(token *ResetToken) (err error) {
	err = pwm.tokencollection.Insert(token)
	return
}

// FindResetToken find reset token by token
func (pwm *Manager) FindResetToken(token string) (tokenobj *ResetToken, err error) {
	tokenobj = &ResetToken{}
	err = pwm.tokencollection.Find(bson.M{"token": token}).One(tokenobj)
	return
}

// DeleteResetToken delete reset token by token
func (pwm *Manager) DeleteResetToken(token string) (err error) {
	_, err = pwm.tokencollection.RemoveAll(bson.M{"token": token})
	return
}
