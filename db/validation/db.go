package validation

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"

	"crypto/rand"
	"math/big"
	"time"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/tools"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoOngoingPhonenumberValidationCollectionName  = "ongoingphonenumbervalidations"
	mongoValidatedPhonenumbers                       = "validatedphonenumbers"
	mongoOngoingEmailAddressValidationCollectionName = "ongoingemailaddressvalidations"
	mongoValidatedEmailAddresses                     = "validatedemailaddresses"
	mongoOngoingAddressValidationCollectionName      = "ongoingaddressvalidations"
	mongoValidatedAddresses                          = "validatedaddresses"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	index := mgo.Index{
		Key:      []string{"key"},
		Unique:   true,
		DropDups: false,
	}

	db.EnsureIndex(mongoOngoingPhonenumberValidationCollectionName, index)
	db.EnsureIndex(mongoOngoingEmailAddressValidationCollectionName, index)
	db.EnsureIndex(mongoOngoingAddressValidationCollectionName, index)

	automaticExpiration := mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 60 * 10,
		Background:  true,
	}
	db.EnsureIndex(mongoOngoingPhonenumberValidationCollectionName, automaticExpiration)

	automaticExpiration = mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 3600 * 48,
		Background:  true,
	}
	db.EnsureIndex(mongoOngoingEmailAddressValidationCollectionName, automaticExpiration)

	automaticExpiration = mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 60 * 60 * 24 * 14, // 14 days
		Background:  true,
	}
	db.EnsureIndex(mongoOngoingAddressValidationCollectionName, automaticExpiration)

	index = mgo.Index{
		Key:      []string{"username", "phonenumber"},
		Unique:   true,
		DropDups: true,
	}
	db.EnsureIndex(mongoValidatedPhonenumbers, index)

	index = mgo.Index{
		Key:      []string{"phonenumber"},
		Unique:   true,
		DropDups: true,
	}
	db.EnsureIndex(mongoValidatedPhonenumbers, index)

	index = mgo.Index{
		Key:      []string{"username", "emailaddress"},
		Unique:   true,
		DropDups: true,
	}
	db.EnsureIndex(mongoValidatedEmailAddresses, index)

	index = mgo.Index{
		Key:      []string{"emailaddress"},
		Unique:   true,
		DropDups: true,
	}
	db.EnsureIndex(mongoValidatedEmailAddresses, index)

	index = mgo.Index{
		Key:      []string{"username", "address"},
		Unique:   true,
		DropDups: true,
	}
	db.EnsureIndex(mongoValidatedAddresses, index)

	// multiple users can have the same physicall address therefore this must not be
	// a unique index
	index = mgo.Index{
		Key:      []string{"address"},
		Unique:   false,
		DropDups: false,
	}
	db.EnsureIndex(mongoValidatedAddresses, index)

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

func (manager *Manager) NewPhonenumberValidationInformation(username string, phonenumber user.Phonenumber) (info *PhonenumberValidationInformation, err error) {
	info = &PhonenumberValidationInformation{CreatedAt: time.Now(), Username: username, Phonenumber: phonenumber.Phonenumber}
	info.Key, err = tools.GenerateRandomString()
	if err != nil {
		return
	}
	numbercode, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return
	}
	info.SMSCode = fmt.Sprintf("%06d", numbercode)
	return
}

func (manager *Manager) NewEmailAddressValidationInformation(username string, email string) (info *EmailAddressValidationInformation, err error) {
	info = &EmailAddressValidationInformation{CreatedAt: time.Now(), Username: username, EmailAddress: email}
	info.Key, err = tools.GenerateRandomString()
	if err != nil {
		return
	}
	info.Secret, err = tools.GenerateRandomString()
	if err != nil {
		return
	}
	return
}

func (manager *Manager) NewAddressValidationInformation(username string, address user.Address) (info *AddressValidationInformation, err error) {
	info = &AddressValidationInformation{CreatedAt: time.Now(), Username: username, Address: address}
	info.Key, err = tools.GenerateRandomString()
	if err != nil {
		return
	}
	info.Secret, err = tools.GenerateRandomString()
	if err != nil {
		return
	}
	return
}

//SavePhonenumberValidationInformation stores a validated phonenumber.
func (manager *Manager) SavePhonenumberValidationInformation(info *PhonenumberValidationInformation) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingPhonenumberValidationCollectionName)
	_, err = mgoCollection.Upsert(bson.M{"phonenumber": info.Phonenumber}, info)
	return
}

//SaveEmailAddressValidationInformation stores a validated emailaddress
func (manager *Manager) SaveEmailAddressValidationInformation(info *EmailAddressValidationInformation) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingEmailAddressValidationCollectionName)
	_, err = mgoCollection.Upsert(bson.M{"emailaddress": info.EmailAddress}, info)
	return
}

//SaveAddressValidationInformation stores a validated address
func (manager *Manager) SaveAddressValidationInformation(info *AddressValidationInformation) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingAddressValidationCollectionName)
	_, err = mgoCollection.Upsert(bson.M{"address": info.Address}, info)
	return
}

func (manager *Manager) RemovePhonenumberValidationInformation(key string) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingPhonenumberValidationCollectionName)
	_, err = mgoCollection.RemoveAll(bson.M{"key": key})
	return
}

func (manager *Manager) RemoveValidatedPhonenumber(username string, phonenumber string) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedPhonenumbers)
	_, err = mgoCollection.RemoveAll(bson.M{"username": username, "phonenumber": phonenumber})
	return
}

func (manager *Manager) RemoveValidatedEmailAddress(username string, email string) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedEmailAddresses)
	_, err = mgoCollection.RemoveAll(bson.M{"username": username, "emailaddress": email})
	return
}

func (manager *Manager) RemoveEmailAddressValidationInformation(key string) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingEmailAddressValidationCollectionName)
	_, err = mgoCollection.RemoveAll(bson.M{"key": key})
	return
}

func (manager *Manager) RemoveValidatedAddress(username string, address user.Address) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedAddresses)
	_, err = mgoCollection.RemoveAll(bson.M{"username": username, "address": address})
	return
}

func (manager *Manager) RemoveAddressValidationInformation(key string) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingAddressValidationCollectionName)
	_, err = mgoCollection.RemoveAll(bson.M{"key": key})
	return
}

func (manager *Manager) UpdatePhonenumberValidationInformation(key string, confirmed bool) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingPhonenumberValidationCollectionName)
	err = mgoCollection.Update(bson.M{"key": key}, bson.M{"$set": bson.M{"confirmed": confirmed}})
	return
}

func (manager *Manager) UpdateEmailAddressValidationInformation(key string, confirmed bool) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingEmailAddressValidationCollectionName)
	err = mgoCollection.Update(bson.M{"key": key}, bson.M{"$set": bson.M{"confirmed": confirmed}})
	return
}

func (manager *Manager) UpdateAddressValidationInformation(key string, confirmed bool) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingAddressValidationCollectionName)
	err = mgoCollection.Update(bson.M{"key": key}, bson.M{"$set": bson.M{"confirmed": confirmed}})
	return
}

func (manager *Manager) GetByKeyPhonenumberValidationInformation(key string) (info *PhonenumberValidationInformation, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingPhonenumberValidationCollectionName)
	err = mgoCollection.Find(bson.M{"key": key}).One(&info)
	if err == mgo.ErrNotFound {
		info = nil
		err = nil
	}
	return
}

func (manager *Manager) GetByKeyEmailAddressValidationInformation(key string) (info *EmailAddressValidationInformation, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingEmailAddressValidationCollectionName)
	err = mgoCollection.Find(bson.M{"key": key}).One(&info)
	if err == mgo.ErrNotFound {
		info = nil
		err = nil
	}
	return
}

func (manager *Manager) GetByKeyAddressValidationInformation(key string) (info *AddressValidationInformation, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingAddressValidationCollectionName)
	err = mgoCollection.Find(bson.M{"key": key}).One(&info)
	if err == mgo.ErrNotFound {
		info = nil
		err = nil
	}
	return
}

func (manager *Manager) GetOngoingEmailAddressValidationByUser(username string) (ongoingValidations []EmailAddressValidationInformation, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingEmailAddressValidationCollectionName)
	err = mgoCollection.Find(bson.M{"username": username}).All(&ongoingValidations)
	return
}

func (manager *Manager) GetOngoingAddressValidationByUser(username string) (ongoingValidations []AddressValidationInformation, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoOngoingAddressValidationCollectionName)
	err = mgoCollection.Find(bson.M{"username": username}).All(&ongoingValidations)
	return
}

func (manager *Manager) NewValidatedPhonenumber(username string, phonenumber string) (validatedphonenumber *ValidatedPhonenumber) {
	validatedphonenumber = &ValidatedPhonenumber{CreatedAt: time.Now(), Username: username, Phonenumber: string(phonenumber)}
	return
}

func (manager *Manager) NewValidatedEmailAddress(username string, email string) (validatedemail *ValidatedEmailAddress) {
	validatedemail = &ValidatedEmailAddress{CreatedAt: time.Now(), Username: username, EmailAddress: email}
	return
}

func (manager *Manager) NewValidatedAddress(username string, address user.Address) (validatedemail *ValidatedAddress) {
	validatedemail = &ValidatedAddress{CreatedAt: time.Now(), Username: username, Address: address}
	return
}

func (manager *Manager) SaveValidatedPhonenumber(validated *ValidatedPhonenumber) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedPhonenumbers)
	_, err = mgoCollection.Upsert(bson.M{"phonenumber": validated.Phonenumber}, validated)
	return
}

func (manager *Manager) SaveValidatedEmailAddress(validated *ValidatedEmailAddress) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedEmailAddresses)
	_, err = mgoCollection.Upsert(bson.M{"emailaddress": validated.EmailAddress}, validated)
	return
}

func (manager *Manager) SaveValidatedAddress(validated *ValidatedAddress) (err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedAddresses)
	_, err = mgoCollection.Upsert(bson.M{"address": validated.Address}, validated)
	return
}

func (manager *Manager) GetByUsernameValidatedPhonenumbers(username string) (validatednumbers []ValidatedPhonenumber, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedPhonenumbers)
	err = mgoCollection.Find(bson.M{"username": username}).All(&validatednumbers)
	return
}

// Returns all validated phone numbers for every username
func (manager *Manager) GetValidatedPhoneNumbersByUsernames(usernames []string) (validatedPhonenumbers []ValidatedPhonenumber, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedPhonenumbers)
	qry := bson.M{
		"username": bson.M{
			"$in": usernames,
		},
	}
	err = mgoCollection.Find(qry).All(&validatedPhonenumbers)
	return
}

func (manager *Manager) GetByEmailAddressValidatedEmailAddress(email string) (validatedemail *ValidatedEmailAddress, err error) {
	validatedemail = &ValidatedEmailAddress{}
	mgoCollection := db.GetCollection(manager.session, mongoValidatedEmailAddresses)
	err = mgoCollection.Find(bson.M{"emailaddress": email}).One(validatedemail)
	return
}

// Returns all validated email addresses for every username
func (manager *Manager) GetValidatedEmailAddressesByUsernames(usernames []string) (validatedEmails []ValidatedEmailAddress, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedEmailAddresses)
	qry := bson.M{
		"username": bson.M{
			"$in": usernames,
		},
	}
	// Avoid returning multiple email addresses for the same user
	err = mgoCollection.Find(qry).All(&validatedEmails)
	return
}

func (manager *Manager) GetByAddressValidatedAddress(address user.Address) (validatedaddress *ValidatedAddress, err error) {
	validatedaddress = &ValidatedAddress{}
	mgoCollection := db.GetCollection(manager.session, mongoValidatedAddresses)
	err = mgoCollection.Find(bson.M{"address": address}).One(validatedaddress)
	return
}

func (manager *Manager) GetByUsernameValidatedEmailAddress(username string) (validatedemails []ValidatedEmailAddress, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedEmailAddresses)
	err = mgoCollection.Find(bson.M{"username": username}).All(&validatedemails)
	return
}

func (manager *Manager) GetByUsernameValidatedAddress(username string) (validatedaddresses []ValidatedAddress, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedAddresses)
	err = mgoCollection.Find(bson.M{"username": username}).All(&validatedaddresses)
	return
}

func (manager *Manager) IsEmailAddressValidated(username string, emailaddress string) (validated bool, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedEmailAddresses)
	count, err := mgoCollection.Find(bson.M{"username": username, "emailaddress": emailaddress}).Count()
	validated = false
	if err != nil {
		return
	}
	if count != 0 {
		validated = true
	}
	return
}

func (manager *Manager) IsAddressValidated(username string, address user.Address) (validated bool, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedAddresses)
	count, err := mgoCollection.Find(bson.M{"username": username, "address": address}).Count()
	validated = false
	if err != nil {
		return
	}
	if count != 0 {
		validated = true
	}
	return
}

func (manager *Manager) IsPhonenumberValidated(username string, phonenumber string) (validated bool, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedPhonenumbers)
	count, err := mgoCollection.Find(bson.M{"username": username, "phonenumber": phonenumber}).Count()
	validated = false
	if err != nil {
		return
	}
	if count != 0 {
		validated = true
	}
	return
}

func (manager *Manager) HasValidatedPhones(username string) (hasValidatedPhones bool, err error) {
	count, err := db.GetCollection(manager.session, mongoValidatedPhonenumbers).Find(bson.M{"username": username}).Count()
	hasValidatedPhones = count > 0
	return hasValidatedPhones, err
}

func (manager *Manager) GetByPhoneNumber(searchString string) (validatedPhonenumber ValidatedPhonenumber, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedPhonenumbers)
	err = mgoCollection.Find(bson.M{"phonenumber": searchString}).One(&validatedPhonenumber)
	return validatedPhonenumber, err
}

func (manager *Manager) GetByEmailAddress(searchString string) (validatedEmailaddress ValidatedEmailAddress, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedEmailAddresses)
	err = mgoCollection.Find(bson.M{"emailaddress": searchString}).One(&validatedEmailaddress)
	return validatedEmailaddress, err
}

func (manager *Manager) GetByAddress(address user.Address) (validatedAddress ValidatedAddress, err error) {
	mgoCollection := db.GetCollection(manager.session, mongoValidatedAddresses)
	err = mgoCollection.Find(bson.M{"address": address}).One(&validatedAddress)
	return validatedAddress, err
}

func (manager *Manager) IsErrNotFound(err error) bool {
	return err == mgo.ErrNotFound
}
