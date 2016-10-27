package user

import (
	"encoding/json"
	"net/http"

	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/communication"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/credentials/totp"
	contractdb "github.com/itsyouonline/identityserver/db/contract"
	organizationDb "github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/registry"
	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/db/user/apikey"
	validationdb "github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/identityservice/contract"
	"github.com/itsyouonline/identityserver/identityservice/invitations"
	"github.com/itsyouonline/identityserver/validation"
	"gopkg.in/mgo.v2"
)

type UsersAPI struct {
	SmsService                    communication.SMSService
	PhonenumberValidationService  *validation.IYOPhonenumberValidationService
	EmailService                  communication.EmailService
	EmailAddressValidationService *validation.IYOEmailAddressValidationService
}

func isUniquePhonenumber(user *user.User, number string, label string) (unique bool) {
	unique = true
	for _, phonenumber := range user.Phonenumbers {
		if phonenumber.Label != label && phonenumber.Phonenumber == number {
			unique = false
			return
		}
	}
	return
}

func isLastVerifiedPhoneNumber(user *user.User, number string, label string, r *http.Request) (last bool, err error) {
	last = false
	valMgr := validationdb.NewManager(r)
	validated, err := valMgr.IsPhonenumberValidated(user.Username, string(number))
	if err != nil {
		return
	}
	if validated {
		// check if this phone number is the last verified one
		uniquelabel := isUniquePhonenumber(user, number, label)
		hasotherverifiednumbers := false
		verifiednumbers, err := valMgr.GetByUsernameValidatedPhonenumbers(user.Username)
		if err != nil {
			return false, err

		}
		for _, verifiednumber := range verifiednumbers {
			if verifiednumber.Phonenumber != string(number) {
				hasotherverifiednumbers = true
				break

			}
		}
		if uniquelabel && !hasotherverifiednumbers {
			return true, nil
		}

	}
	return
}

// It is handler for POST /users
func (api UsersAPI) Post(w http.ResponseWriter, r *http.Request) {

	var u user.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userMgr := user.NewManager(r)
	if err := userMgr.Save(&u); err != nil {
		log.Error("ERROR while saving user:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(&u)
}

// GetUser is handler for GET /users/{username}
func (api UsersAPI) GetUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	userMgr := user.NewManager(r)

	user, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func isValidLabel(label string) (valid bool) {
	valid = true
	labelLength := len(label)
	valid = valid && labelLength > 1 && labelLength < 51

	if !valid {
		log.Debug("Invalid label: ", label)
	}
	return valid
}

// RegisterNewEmailAddress is the handler for POST /users/{username}/emailaddresses
// Register a new email address
func (api UsersAPI) RegisterNewEmailAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	body := user.EmailAddress{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(body.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userMgr := user.NewManager(r)
	u, err := userMgr.GetByName(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := u.GetEmailAddressByLabel(body.Label); err == nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if err = userMgr.SaveEmail(username, body); err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	valMgr := validationdb.NewManager(r)
	validated, err := valMgr.IsEmailAddressValidated(username, body.EmailAddress)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !validated {
		_, err = api.EmailAddressValidationService.RequestValidation(r, username, body.EmailAddress, fmt.Sprintf("https://%s/emailvalidation", r.Host))
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

// UpdateEmailAddress is the handler for PUT /users/{username}/emailaddresses/{label}
// Updates the label and/or value of an email address
func (api UsersAPI) UpdateEmailAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	oldlabel := mux.Vars(r)["label"]

	body := user.EmailAddress{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(body.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userMgr := user.NewManager(r)

	if oldlabel != body.Label {
		u, err := userMgr.GetByName(username)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if _, err := u.GetEmailAddressByLabel(body.Label); err == nil {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	if err := userMgr.SaveEmail(username, body); err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if oldlabel != body.Label {
		if err := userMgr.RemoveEmail(username, oldlabel); err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	valMgr := validationdb.NewManager(r)
	validated, err := valMgr.IsEmailAddressValidated(username, body.EmailAddress)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !validated {
		_, err = api.EmailAddressValidationService.RequestValidation(r, username, body.EmailAddress, fmt.Sprintf("https://%s/emailvalidation", r.Host))
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

// Validate email address is the handler for GET /users/{username}/emailaddress/{label}/validate
func (api UsersAPI) ValidateEmailAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	email, err := userobj.GetEmailAddressByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err = api.EmailAddressValidationService.RequestValidation(r, username, email.EmailAddress, fmt.Sprintf("https://%s/emailvalidation", r.Host))
	w.WriteHeader(http.StatusNoContent)
}

// ListEmailAddresses is the handler for GET /users/{username}/emailaddresses
func (api UsersAPI) ListEmailAddresses(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	validated := strings.Contains(r.URL.RawQuery, "validated")
	userMgr := user.NewManager(r)
	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	var emails []user.EmailAddress
	if validated {
		emails = make([]user.EmailAddress, 0)
		valMngr := validationdb.NewManager(r)
		validatedemails, err := valMngr.GetByUsernameValidatedEmailAddress(username)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		for _, email := range userobj.EmailAddresses {
			for _, validatedemail := range validatedemails {
				if email.EmailAddress == validatedemail.EmailAddress {
					emails = append(emails, email)
					break
				}
			}
		}
	} else {
		emails = userobj.EmailAddresses
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emails)
}

// DeleteEmailAddress is the handler for DELETE /users/{username}/emailaddresses/{label}
// Removes an email address
func (api UsersAPI) DeleteEmailAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]

	userMgr := user.NewManager(r)
	valMgr := validationdb.NewManager(r)

	u, err := userMgr.GetByName(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	email, err := u.GetEmailAddressByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if len(u.EmailAddresses) == 1 {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if err = userMgr.RemoveEmail(username, label); err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err = valMgr.RemoveValidatedEmailAddress(username, email.EmailAddress); err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// DeleteGithubAccount is the handler for DELETE /users/{username}/github
// Delete the associated Github account.
func (api UsersAPI) DeleteGithubAccount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)
	err := userMgr.DeleteGithubAccount(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteFacebookAccount is the handler for DELETE /users/{username}/facebook
// Delete the associated facebook account
func (api UsersAPI) DeleteFacebookAccount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	userMgr := user.NewManager(r)
	err := userMgr.DeleteFacebookAccount(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePassword is the handler for PUT /users/{username}/password
func (api UsersAPI) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	body := struct {
		Currentpassword string `json:"currentpassword"`
		Newpassword     string `json:"newpassword"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	userMgr := user.NewManager(r)
	exists, err := userMgr.Exists(username)
	if !exists || err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	passwordMgr := password.NewManager(r)
	passwordok, err := passwordMgr.Validate(username, body.Currentpassword)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !passwordok {
		writeErrorResponse(w, 422, "incorrect_password")
		return
	}
	err = passwordMgr.Save(username, body.Newpassword)
	if err != nil {
		writeErrorResponse(w, 422, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetUserInformation is the handler for GET /users/{username}/info
func (api UsersAPI) GetUserInformation(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)

	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	requestingClient, validClient := context.Get(r, "client_id").(string)
	if !validClient {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	authorization, err := userMgr.GetAuthorization(username, requestingClient)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respBody := &Userview{
		Username: userobj.Username,
	}

	if authorization.Name {
		respBody.Firstname = userobj.Firstname
		respBody.Lastname = userobj.Lastname
	}
	if authorization.Github {
		respBody.Github = userobj.Github
	}
	if authorization.Facebook {
		respBody.Facebook = userobj.Facebook
	}
	if authorization.Addresses != nil {
		respBody.Addresses = make([]user.Address, 0)

		for _, addressmap := range authorization.Addresses {
			address, err := userobj.GetAddressByLabel(addressmap.RealLabel)
			if err == nil {
				newaddress := user.Address{
					Label:      addressmap.RequestedLabel,
					City:       address.City,
					Country:    address.Country,
					Nr:         address.Nr,
					Other:      address.Other,
					Postalcode: address.Postalcode,
					Street:     address.Street,
				}
				respBody.Addresses = append(respBody.Addresses, newaddress)
			} else {
				log.Debug(err)
			}
		}
	}

	if authorization.EmailAddresses != nil {
		respBody.EmailAddresses = make([]user.EmailAddress, 0)

		for _, emailmap := range authorization.EmailAddresses {
			email, err := userobj.GetEmailAddressByLabel(emailmap.RealLabel)
			if err == nil {
				newemail := user.EmailAddress{
					Label:        emailmap.RequestedLabel,
					EmailAddress: email.EmailAddress,
				}
				respBody.EmailAddresses = append(respBody.EmailAddresses, newemail)
			} else {
				log.Debug(err)
			}
		}
	}

	if authorization.Phonenumbers != nil {
		respBody.Phonenumbers = make([]user.Phonenumber, 0)
		for _, phonemap := range authorization.Phonenumbers {
			phonenumber, err := userobj.GetPhonenumberByLabel(phonemap.RealLabel)
			if err == nil {
				newnumber := user.Phonenumber{
					Label:       phonemap.RequestedLabel,
					Phonenumber: phonenumber.Phonenumber,
				}
				respBody.Phonenumbers = append(respBody.Phonenumbers, newnumber)
			} else {
				log.Debug(err)
			}
		}
	}

	if authorization.BankAccounts != nil {
		respBody.BankAccounts = make([]user.BankAccount, 0)

		for _, bankmap := range authorization.BankAccounts {
			bank, err := userobj.GetBankAccountByLabel(bankmap.RealLabel)
			if err == nil {
				newbank := user.BankAccount{
					Label:   bankmap.RequestedLabel,
					Bic:     bank.Bic,
					Country: bank.Country,
					Iban:    bank.Iban,
				}
				respBody.BankAccounts = append(respBody.BankAccounts, newbank)
			} else {
				log.Debug(err)
			}
		}
	}
	if authorization.DigitalWallet != nil {
		respBody.DigitalWallet = make([]user.DigitalAssetAddress, 0)

		for _, addressMap := range authorization.DigitalWallet {
			walletAddress, err := userobj.GetDigitalAssetAddressByLabel(addressMap.RealLabel)
			if err == nil {
				walletAddress.Label = addressMap.RequestedLabel
				respBody.DigitalWallet = append(respBody.DigitalWallet, walletAddress)
			} else {
				log.Debug(err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respBody)
}

// ValidateUsername is the handler for GET /users/{username}/validate
func (api UsersAPI) ValidateUsername(w http.ResponseWriter, r *http.Request) {

	// token := req.FormValue("token")

	// uncomment below line to add header
	// w.Header.Set("key","value")
}

// RegisterNewPhonenumber is the handler for POST /users/{username}/phonenumbers
// Register a new phonenumber
func (api UsersAPI) RegisterNewPhonenumber(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)

	u, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	body := user.Phonenumber{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(body.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !body.IsValid() {
		log.Debug("Invalid phonenumber: ", body.Phonenumber)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Check if this label is already used
	_, err = u.GetPhonenumberByLabel(body.Label)
	if err == nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if err := userMgr.SavePhone(username, body); err != nil {
		log.Error("ERROR while saving a phonenumber - ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// respond with created phone number.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(body)
}

// GetUserPhoneNumbers is the handler for GET /users/{username}/phonenumbers
func (api UsersAPI) GetUserPhoneNumbers(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	validated := strings.Contains(r.URL.RawQuery, "validated")
	userMgr := user.NewManager(r)

	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	var phonenumbers []user.Phonenumber
	if validated {
		phonenumbers = make([]user.Phonenumber, 0)
		valMngr := validationdb.NewManager(r)
		validatednumbers, err := valMngr.GetByUsernameValidatedPhonenumbers(username)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		for _, number := range userobj.Phonenumbers {
			for _, validatednumber := range validatednumbers {
				if number.Phonenumber == validatednumber.Phonenumber {
					phonenumbers = append(phonenumbers, number)
					break
				}
			}
		}
	} else {
		phonenumbers = userobj.Phonenumbers
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(phonenumbers)
}

// GetUserPhonenumberByLabel is the handler for GET /users/{username}/phonenumbers/{label}
func (api UsersAPI) GetUserPhonenumberByLabel(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	phonenumber, err := userobj.GetPhonenumberByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(phonenumber)
}

// ValidatePhoneNumber is the handler for POST /users/{username}/phonenumbers/{label}/validate
func (api UsersAPI) ValidatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	phonenumber, err := userobj.GetPhonenumberByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	validationKey := ""
	validationKey, err = api.PhonenumberValidationService.RequestValidation(r, username, phonenumber, fmt.Sprintf("https://%s/phonevalidation", r.Host))
	response := struct {
		ValidationKey string `json:"validationkey"`
	}{
		ValidationKey: validationKey,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}

// VerifyPhoneNumber is the handler for PUT /users/{username}/phonenumbers/{label}/validate
func (api UsersAPI) VerifyPhoneNumber(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	values := struct {
		Smscode       string `json:"smscode"`
		ValidationKey string `json:"validationkey"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ProcessPhonenumberConfirmation request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err = userobj.GetPhonenumberByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err = api.PhonenumberValidationService.ConfirmValidation(r, values.ValidationKey, values.Smscode)
	if err != nil {
		log.Debug(err)
		if err == validation.ErrInvalidCode || err == validation.ErrInvalidOrExpiredKey {
			writeErrorResponse(w, 422, err.Error())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdatePhonenumber is the handler for PUT /users/{username}/phonenumbers/{label}
// Update the label and/or value of an existing phonenumber.
func (api UsersAPI) UpdatePhonenumber(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	oldlabel := mux.Vars(r)["label"]

	body := user.Phonenumber{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(body.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !body.IsValid() {
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	userMgr := user.NewManager(r)

	u, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	oldnumber, err := u.GetPhonenumberByLabel(oldlabel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if oldlabel != body.Label {
		// Check if there already is another phone number with the new label
		_, err := u.GetPhonenumberByLabel(body.Label)
		if err == nil {
			writeErrorResponse(w, http.StatusConflict, "duplicate_label")
			return
		}
	}

	if oldnumber.Phonenumber != body.Phonenumber {
		last, err := isLastVerifiedPhoneNumber(u, oldnumber.Phonenumber, oldlabel, r)
		if err != nil {
			log.Error("ERROR while verifying last verified number - ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if last {
			writeErrorResponse(w, http.StatusConflict, "cannot_delete_last_verified_phone_number")
			return
		}
	}

	if err = userMgr.SavePhone(username, body); err != nil {
		log.Error("ERROR while saving phonenumber - ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if oldlabel != body.Label {
		if err := userMgr.RemovePhone(username, oldlabel); err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	valMgr := validationdb.NewManager(r)
	if oldnumber.Phonenumber != body.Phonenumber && isUniquePhonenumber(u, oldnumber.Phonenumber, oldlabel) {
		valMgr.RemoveValidatedPhonenumber(username, oldnumber.Phonenumber)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)

}

// DeletePhonenumber is the handler for DELETE /users/{username}/phonenumbers/{label}
// Removes a phonenumber
func (api UsersAPI) DeletePhonenumber(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)
	valMgr := validationdb.NewManager(r)
	force := r.URL.Query().Get("force") == "true"

	user, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	number, err := user.GetPhonenumberByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	last, err := isLastVerifiedPhoneNumber(user, number.Phonenumber, label, r)
	if err != nil {
		log.Error("ERROR while checking if number can be deleted:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
	if last {
		hasTOTP := false
		if !force {
			writeErrorResponse(w, http.StatusConflict, "warning_delete_last_verified_phone_number")
			return
		} else {
			totpMgr := totp.NewManager(r)
			hasTOTP, err = totpMgr.HasTOTP(username)
		}
		if !hasTOTP {
			writeErrorResponse(w, http.StatusConflict, "cannot_delete_last_verified_phone_number")
			return
		}
	}

	// check if the phonenumber is unique or if there are duplicates
	uniqueNumber := isUniquePhonenumber(user, number.Phonenumber, label)

	if err := userMgr.RemovePhone(username, label); err != nil {
		log.Error("ERROR while saving user:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// only remove the phonenumber from the validatedphonenumber collection if there are no duplicates
	if uniqueNumber {
		if err := valMgr.RemoveValidatedPhonenumber(username, number.Phonenumber); err != nil {
			log.Error("ERROR while saving user:\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateUserBankAccount is handler for POST /users/{username}/banks
// Create new bank account
func (api UsersAPI) CreateUserBankAccount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)

	bank := user.BankAccount{}

	if err := json.NewDecoder(r.Body).Decode(&bank); err != nil {
		log.Error("Error while decoding the body: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(bank.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	//Check if this label is already used
	_, err = user.GetBankAccountByLabel(bank.Label)
	if err == nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if err := userMgr.SaveBank(user, bank); err != nil {
		log.Error("ERROR while saving address:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// respond with created bank account
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(bank)
}

// GetUserBankAccounts It is handler for GET /users/{username}/banks
func (api UsersAPI) GetUserBankAccounts(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)

	user, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.BankAccounts)
}

// GetUserBankAccountByLabel is handler for GET /users/{username}/banks/{label}
func (api UsersAPI) GetUserBankAccountByLabel(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	bank, err := userobj.GetBankAccountByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bank)
}

// UpdateUserBankAccount is handler for PUT /users/{username}/banks/{label}
// Update an existing bankaccount and label.
func (api UsersAPI) UpdateUserBankAccount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	oldlabel := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	newbank := user.BankAccount{}

	if err := json.NewDecoder(r.Body).Decode(&newbank); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(newbank.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	oldbank, err := user.GetBankAccountByLabel(oldlabel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if oldbank.Label != newbank.Label {
		_, err := user.GetBankAccountByLabel(newbank.Label)
		if err == nil {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	if err = userMgr.SaveBank(user, newbank); err != nil {
		log.Error("ERROR while saving bank - ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if oldlabel != newbank.Label {
		if err := userMgr.RemoveBank(user, oldlabel); err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newbank)
}

// DeleteUserBankAccount is handler for DELETE /users/{username}/banks/{label}
// Delete a BankAccount
func (api UsersAPI) DeleteUserBankAccount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	user, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err = user.GetBankAccountByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := userMgr.RemoveBank(user, label); err != nil {
		log.Error("ERROR while saving user:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterNewAddress is the handler for POST /users/{username}/addresses
// Register a new address
func (api UsersAPI) RegisterNewAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)

	address := user.Address{}

	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		log.Debug("Error while decoding the body: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(address.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	u, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	//Check if this label is already used
	_, err = u.GetAddressByLabel(address.Label)
	if err == nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if err := userMgr.SaveAddress(username, address); err != nil {
		log.Error("ERROR while saving address:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// respond with created phone number.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(address)
}

// GetUserAddresses is handler for GET /users/{username}/addresses
func (api UsersAPI) GetUserAddresses(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)

	user, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.Addresses)
}

// GetUserAddressByLabel is handler for GET /users/{username}/addresses/{label}
func (api UsersAPI) GetUserAddressByLabel(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	address, err := userobj.GetAddressByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

// UpdateAddress is the handler for PUT /users/{username}/addresses/{label}
// Update the label and/or value of an existing address.
func (api UsersAPI) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	oldlabel := mux.Vars(r)["label"]

	newaddress := user.Address{}
	if err := json.NewDecoder(r.Body).Decode(&newaddress); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(newaddress.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userMgr := user.NewManager(r)

	u, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err = u.GetAddressByLabel(oldlabel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if oldlabel != newaddress.Label {
		_, err = u.GetAddressByLabel(newaddress.Label)
		if err == nil {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	if err = userMgr.SaveAddress(username, newaddress); err != nil {
		log.Error("ERROR while saving address - ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if oldlabel != newaddress.Label {
		if err := userMgr.RemoveAddress(username, oldlabel); err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newaddress)
}

// DeleteAddress is the handler for DELETE /users/{username}/addresses/{label}
// Removes an address
func (api UsersAPI) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	u, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err = u.GetAddressByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err = userMgr.RemoveAddress(username, label); err != nil {
		log.Error("ERROR while saving address:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserContracts is handler for GET /users/{username}/contracts
// Get the contracts where the user is 1 of the parties. Order descending by date.
func (api UsersAPI) GetUserContracts(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	includedparty := contractdb.Party{Type: "user", Name: username}
	contract.FindContracts(w, r, includedparty)
}

// RegisterNewContract is handler for GET /users/{username}/contracts
func (api UsersAPI) RegisterNewContract(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	includedparty := contractdb.Party{Type: "user", Name: username}
	contract.CreateContract(w, r, includedparty)

}

// GetNotifications is handler for GET /users/{username}/notifications
// Get the list of notifications, these are pending invitations or approvals
func (api UsersAPI) GetNotifications(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	type NotificationList struct {
		Approvals        []invitations.JoinOrganizationInvitation `json:"approvals"`
		ContractRequests []contractdb.ContractSigningRequest      `json:"contractRequests"`
		Invitations      []invitations.JoinOrganizationInvitation `json:"invitations"`
	}
	var notifications NotificationList

	invititationMgr := invitations.NewInvitationManager(r)

	userOrgRequests, err := invititationMgr.GetByUser(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	notifications.Invitations = userOrgRequests

	// TODO: Get Approvals and Contract requests
	notifications.Approvals = []invitations.JoinOrganizationInvitation{}
	notifications.ContractRequests = []contractdb.ContractSigningRequest{}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(&notifications)

}

// usernameorganizationsGet is the handler for GET /users/{username}/organizations
// Get the list organizations a user is owner of member of
func (api UsersAPI) usernameorganizationsGet(w http.ResponseWriter, r *http.Request) {

}

// GetAllAuthorizations is the handler for GET /users/{username}/authorizations
// Get the list of authorizations.
func (api UsersAPI) GetAllAuthorizations(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	userMgr := user.NewManager(r)

	authorizations, err := userMgr.GetAuthorizationsByUser(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(authorizations)

}

// GetAuthorization is the handler for GET /users/{username}/authorizations/{grantedTo}
// Get the authorization for a specific organization.
func (api UsersAPI) GetAuthorization(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	grantedTo := mux.Vars(r)["grantedTo"]

	userMgr := user.NewManager(r)

	authorization, err := userMgr.GetAuthorization(username, grantedTo)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if authorization == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(authorization)
}

// UpdateAuthorization is the handler for PUT /users/{username}/authorizations/{grantedTo}
// Modify which information an organization is able to see.
func (api UsersAPI) UpdateAuthorization(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	grantedTo := mux.Vars(r)["grantedTo"]

	authorization := &user.Authorization{}

	if err := json.NewDecoder(r.Body).Decode(authorization); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	authorization.Username = username
	authorization.GrantedTo = grantedTo

	userMgr := user.NewManager(r)

	err := userMgr.UpdateAuthorization(authorization)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authorization)
}

// DeleteAuthorization is the handler for DELETE /users/{username}/authorizations/{grantedTo}
// Remove the authorization for an organization, the granted organization will no longer
// have access the user's information.
func (api UsersAPI) DeleteAuthorization(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	grantedTo := mux.Vars(r)["grantedTo"]

	userMgr := user.NewManager(r)

	err := userMgr.DeleteAuthorization(username, grantedTo)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (api UsersAPI) AddAPIKey(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	body := struct {
		Label string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	apikeyMgr := apikey.NewManager(r)
	apikey := apikey.NewAPIKey(username, body.Label)
	apikeyMgr.Save(apikey)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apikey)
}

func (api UsersAPI) GetAPIKey(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	apikeyMgr := apikey.NewManager(r)
	apikey, err := apikeyMgr.GetByUsernameAndLabel(username, label)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apikey)
}

func (api UsersAPI) UpdateAPIKey(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	apikeyMgr := apikey.NewManager(r)
	body := struct {
		Label string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	apikey, err := apikeyMgr.GetByUsernameAndLabel(username, label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	apikey.Label = body.Label
	apikeyMgr.Save(apikey)
	w.WriteHeader(http.StatusNoContent)

}
func (api UsersAPI) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	apikeyMgr := apikey.NewManager(r)
	apikeyMgr.Delete(username, label)

}
func (api UsersAPI) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	apikeyMgr := apikey.NewManager(r)
	apikeys, err := apikeyMgr.GetByUser(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if apikeys == nil {
		apikeys = []apikey.APIKey{}
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(apikeys)
}

// UpdateName is the handler for PUT /users/{username}/name
func (api UsersAPI) UpdateName(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	values := struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	userMgr := user.NewManager(r)
	exists, err := userMgr.Exists(username)
	if !exists || err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	err = userMgr.UpdateName(username, values.Firstname, values.Lastname)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetTwoFAMethods is the handler for GET /users/{username}/twofamethods
// Get the possible two factor authentication methods
func (api UsersAPI) GetTwoFAMethods(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)
	userFromDB, err := userMgr.GetByName(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := struct {
		Totp bool               `json:"totp"`
		Sms  []user.Phonenumber `json:"sms"`
	}{}
	totpMgr := totp.NewManager(r)
	response.Totp, err = totpMgr.HasTOTP(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	valMgr := validationdb.NewManager(r)
	verifiedPhones, err := valMgr.GetByUsernameValidatedPhonenumbers(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for _, validatedPhoneNumber := range verifiedPhones {
		for _, number := range userFromDB.Phonenumbers {
			if number.Phonenumber == string(validatedPhoneNumber.Phonenumber) {
				response.Sms = append(response.Sms, number)
			}
		}
	}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
	return
}

// GetTOTPSecret is the handler for GET /users/{username}/totp/
// Gets a new TOTP secret
func (api UsersAPI) GetTOTPSecret(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Totpsecret string `json:"totpsecret"`
	}{}
	token, err := totp.NewToken()
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response.Totpsecret = token.Secret
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}

// SetupTOTP is the handler for POST /users/{username}/totp/
// Configures TOTP authentication for this user
func (api UsersAPI) SetupTOTP(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	values := struct {
		TotpSecret string `json:"totpsecret"`
		TotpCode   string `json:"totpcode"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	totpMgr := totp.NewManager(r)
	err := totpMgr.Save(username, values.TotpSecret)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	valid, err := totpMgr.Validate(username, values.TotpCode)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !valid {
		err := totpMgr.Remove(username)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(422)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// RemoveTOTP is the handler for DELETE /users/{username}/totp/
// Removes TOTP authentication for this user, if possible.
func (api UsersAPI) RemoveTOTP(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	valMngr := validationdb.NewManager(r)
	hasValidatedPhones, err := valMngr.HasValidatedPhones(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !hasValidatedPhones {
		w.WriteHeader(http.StatusConflict)
		return
	}
	totpMgr := totp.NewManager(r)
	err = totpMgr.Remove(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// LeaveOrganization is the handler for DELETE /users/{username}/organizations/{globalid}/leave
// Removes the user from an organization
func (api UsersAPI) LeaveOrganization(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	organizationGlobalid := mux.Vars(r)["globalid"]
	orgMgr := organizationDb.NewManager(r)
	err := orgMgr.RemoveUser(organizationGlobalid, username)
	if err == mgo.ErrNotFound {
		writeErrorResponse(w, http.StatusNotFound, "user_not_found")
		return
	} else if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	userMgr := user.NewManager(r)
	err = userMgr.DeleteAuthorization(username, organizationGlobalid)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListUserRegistry is the handler for GET /users/{username}/registry
// Lists the Registry entries
func (api UsersAPI) ListUserRegistry(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	mgr := registry.NewManager(r)
	registryEntries, err := mgr.ListRegistryEntries(username, "")
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registryEntries)
}

// AddUserRegistryEntry is the handler for POST /users/{username}/registry
// Adds a RegistryEntry to the user's registry, if the key is already used, it is overwritten.
func (api UsersAPI) AddUserRegistryEntry(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	registryEntry := registry.RegistryEntry{}

	if err := json.NewDecoder(r.Body).Decode(&registryEntry); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := registryEntry.Validate(); err != nil {
		log.Debug("Invalid registry entry: ", registryEntry)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mgr := registry.NewManager(r)
	err := mgr.UpsertRegistryEntry(username, "", registryEntry)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registryEntry)
}

// GetUserRegistryEntry is the handler for GET /users/{username}/registry/{key}
// Get a RegistryEntry from the user's registry.
func (api UsersAPI) GetUserRegistryEntry(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	key := mux.Vars(r)["key"]

	mgr := registry.NewManager(r)
	registryEntry, err := mgr.GetRegistryEntry(username, "", key)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registryEntry == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registryEntry)
}

// DeleteUserRegistryEntry is the handler for DELETE /users/{username}/registry/{key}
// Removes a RegistryEntry from the user's registry
func (api UsersAPI) DeleteUserRegistryEntry(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	key := mux.Vars(r)["key"]

	mgr := registry.NewManager(r)
	err := mgr.DeleteRegistryEntry(username, "", key)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeErrorResponse(responseWrite http.ResponseWriter, httpStatusCode int, message string) {
	log.Debug(httpStatusCode, message)
	errorResponse := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}
	responseWrite.WriteHeader(httpStatusCode)
	json.NewEncoder(responseWrite).Encode(&errorResponse)
}
