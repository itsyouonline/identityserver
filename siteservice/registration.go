package siteservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/user"
	validationdb "github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"
	"github.com/itsyouonline/identityserver/validation"
)

const (
	mongoRegistrationCollectionName = "registrationsessions"
	MAX_PENDING_REGISTRATION_COUNT  = 10000
)

//initLoginModels initialize models in mongo
func (service *Service) initRegistrationModels() {
	index := mgo.Index{
		Key:      []string{"sessionkey"},
		Unique:   true,
		DropDups: false,
	}

	db.EnsureIndex(mongoRegistrationCollectionName, index)

	automaticExpiration := mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 60 * 10, //10 minutes
		Background:  true,
	}
	db.EnsureIndex(mongoRegistrationCollectionName, automaticExpiration)

}

type registrationSessionInformation struct {
	SessionKey           string
	SMSCode              string
	Confirmed            bool
	ConfirmationAttempts uint
	CreatedAt            time.Time
}

const (
	registrationFileName = "registration.html"
)

func (service *Service) renderRegistrationFrom(w http.ResponseWriter, request *http.Request) {
	htmlData, err := html.Asset(registrationFileName)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sessions.Save(request, w)
	w.Write(htmlData)
}

//CheckRegistrationSMSConfirmation is called by the sms code form to check if the sms is already confirmed on the mobile phone
func (service *Service) CheckRegistrationSMSConfirmation(w http.ResponseWriter, request *http.Request) {
	registrationSession, err := service.GetSession(request, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := map[string]bool{}

	if registrationSession.IsNew {
		// TODO: registrationSession is new with SMS, something must be wrong
		log.Warn("Registration is new")
		response["confirmed"] = true //This way the form will be submitted, let the form handler deal with redirect to login
	} else {
		validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)

		confirmed, err := service.phonenumberValidationService.IsConfirmed(request, validationkey)
		if err == validation.ErrInvalidOrExpiredKey {
			confirmed = true //This way the form will be submitted, let the form handler deal with redirect to login
			return
		}
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response["confirmed"] = confirmed
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//CheckRegistrationEmailConfirmation is called by the regisration form to check if the email is already confirmed
func (service *Service) CheckRegistrationEmailConfirmation(w http.ResponseWriter, r *http.Request) {
	registrationSession, err := service.GetSession(r, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := map[string]bool{}

	if registrationSession.IsNew {
		// TODO: registrationSession is new, something must be wrong
		log.Warn("Registration is new")
		response["confirmed"] = true //This way the form will be submitted, let the form handler deal with redirect to login
	} else {
		validationkey, _ := registrationSession.Values["emailvalidationkey"].(string)

		confirmed, err := service.emailaddressValidationService.IsConfirmed(r, validationkey)
		if err == validation.ErrInvalidOrExpiredKey {
			// TODO
			confirmed = true //This way the form will be submitted, let the form handler deal with redirect to login
			return
		}
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response["confirmed"] = confirmed
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//ShowRegistrationForm shows the user registration page
func (service *Service) ShowRegistrationForm(w http.ResponseWriter, request *http.Request) {
	service.renderRegistrationFrom(w, request)
}

//ProcessPhonenumberConfirmationForm processes the Phone number confirmation form
func (service *Service) ProcessPhonenumberConfirmationForm(w http.ResponseWriter, r *http.Request) {
	values := struct {
		Smscode string `json:"smscode"`
	}{}

	response := struct {
		Error     string `json:"error"`
		Confirmed bool   `json:"confirmed"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ProcessPhonenumberConfirmation request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	registrationSession, err := service.GetSession(r, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registrationSession.IsNew {
		sessions.Save(r, w)
		log.Debug("Registration session expired")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username, _ := registrationSession.Values["username"].(string)
	validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)

	if isConfirmed, _ := service.phonenumberValidationService.IsConfirmed(r, validationkey); isConfirmed {
		userMgr := user.NewManager(r)
		userMgr.RemoveExpireDate(username)
		response.Confirmed = true
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	smscode := values.Smscode
	if err != nil || smscode == "" {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = service.phonenumberValidationService.ConfirmValidation(r, validationkey, smscode)
	if err == validation.ErrInvalidCode {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response.Error = "invalid_sms_code"
		json.NewEncoder(w).Encode(&response)
		return
	}
	if err == validation.ErrInvalidOrExpiredKey {
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&response)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	userMgr := user.NewManager(r)
	userMgr.RemoveExpireDate(username)
	response.Confirmed = true
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

//ResendPhonenumberConfirmation resend the phonenumberconfirmation to a possbily new phonenumber
func (service *Service) ResendPhonenumberConfirmation(w http.ResponseWriter, request *http.Request) {
	values := struct {
		PhoneNumber string `json:"phonenumber"`
		LangKey     string `json:"langkey"`
	}{}

	response := struct {
		RedirectUrL string `json:"redirecturl"`
		Error       string `json:"error"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ResendPhonenumberConfirmation request: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	registrationSession, err := service.GetSession(request, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registrationSession.IsNew {
		sessions.Save(request, w)
		log.Debug("Registration session expired")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username, _ := registrationSession.Values["username"].(string)

	//Invalidate the previous validation request, ignore a possible error
	validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)
	_ = service.phonenumberValidationService.ExpireValidation(request, validationkey)

	phonenumber := user.Phonenumber{Label: "main", Phonenumber: values.PhoneNumber}
	if !phonenumber.Validate() {
		log.Debug("Invalid phone number")
		w.WriteHeader(http.StatusUnprocessableEntity)
		response.Error = "invalid_phonenumber"
		json.NewEncoder(w).Encode(&response)
		return
	}

	uMgr := user.NewManager(request)
	err = uMgr.SavePhone(username, phonenumber)
	if err != nil {
		log.Error("ResendPhonenumberConfirmation: Could not save phonenumber: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	validationkey, err = service.phonenumberValidationService.RequestValidation(request, username, phonenumber, fmt.Sprintf("https://%s/phonevalidation", request.Host), values.LangKey)
	if err != nil {
		log.Error("ResendPhonenumberConfirmation: Could not get validationkey: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	registrationSession.Values["phonenumbervalidationkey"] = validationkey

	sessions.Save(request, w)
	response.RedirectUrL = fmt.Sprintf("https://%s/register/#smsconfirmation", request.Host)
	json.NewEncoder(w).Encode(&response)
}

//ProcessRegistrationForm processes the user registration form
func (service *Service) ProcessRegistrationForm(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Redirecturl string `json:"redirecturl"`
		Error       string `json:"error"`
	}{}
	values := struct {
		Firstname       string `json:"firstname"`
		Lastname        string `json:"lastname"`
		Email           string `json:"email"`
		Phonenumber     string `json:"phonenumber"`
		PhonenumberCode string `json:"phonenumbercode"`
		Password        string `json:"password"`
		RedirectParams  string `json:"redirectparams"`
		LangKey         string `json:"langkey"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the registration request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	registrationSession, err := service.GetSession(r, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error("Failed to retrieve registration session: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registrationSession.IsNew {
		sessions.Save(r, w)
		log.Debug("Registration session expired")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username, _ := registrationSession.Values["username"].(string)

	userMgr := user.NewManager(r)

	// check if phone number is validated or sms code is provided to validate phone
	phonevalidationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)

	if isConfirmed, _ := service.phonenumberValidationService.IsConfirmed(r, phonevalidationkey); !isConfirmed {

		smscode := values.PhonenumberCode
		if smscode == "" {
			log.Debug("no sms code provided and phone not confirmed yet")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = service.phonenumberValidationService.ConfirmValidation(r, phonevalidationkey, smscode)
		if err == validation.ErrInvalidCode {
			w.WriteHeader(http.StatusUnprocessableEntity)
			response.Error = "invalid_sms_code"
			json.NewEncoder(w).Encode(&response)
			return
		}
		if err == validation.ErrInvalidOrExpiredKey {
			sessions.Save(r, w)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&response)
			return
		}
		if err != nil {
			log.Error("Error while trying to validate phone number in regsitration flow: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	// at this point the phone number is confirmed
	userMgr.RemoveExpireDate(username)
	// Check if the email has already been verified through the link

	emailvalidationkey, _ := registrationSession.Values["emailvalidationkey"].(string)
	emailConfirmed, _ := service.emailaddressValidationService.IsConfirmed(r, emailvalidationkey)
	if !emailConfirmed {
		log.Debug("Email not confirmed yet")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Ideally, we would remove the registration session here as registration is completed.
	// However the login handler checks the existence of this session because it needs the
	// redirectparams as part of the logic to move the user to the requested authenticated page.
	// But this means that if the user immediatly goes back to the registration screen, the old
	// user data is modified as there is already data in the session such as a username. Since we can't
	// remove the session, just empty out al the keys to mimic this process, and then only set the
	// redirectparams

	// Clear registration session
	for key := range registrationSession.Values {
		delete(registrationSession.Values, key)
	}

	// Now set the redirectparams
	registrationSession.Values["redirectparams"] = values.RedirectParams

	sessions.Save(r, w)
	service.loginUser(w, r, username)
}

//ValidateUsername checks if a username is already taken or not
func (service *Service) ValidateUsername(w http.ResponseWriter, request *http.Request) {
	username := request.URL.Query().Get("username")
	response := struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
	}{
		Valid: true,
		Error: "",
	}
	valid := user.ValidateUsername(username)
	if !valid {
		log.Debug("Invalid username format:", username)
		response.Error = "invalid_username_format"
		response.Valid = false
		json.NewEncoder(w).Encode(&response)
		return
	}
	userMgr := user.NewManager(request)
	orgMgr := organization.NewManager(request)
	userExists, err := userMgr.Exists(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if userExists {
		log.Debug("username ", username, " already taken")
		response.Error = "user_exists"
		response.Valid = false
	} else {
		orgExists := orgMgr.Exists(username)
		if orgExists {
			log.Debugf("Organization with name %s already exists", username)
			response.Error = "organization_exists"
			response.Valid = false
		}
	}
	json.NewEncoder(w).Encode(&response)
	return
}

// ValidateInfo starts validation for a temporary username
func (service *Service) ValidateInfo(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Password  string `json:"password"`
		LangKey   string `json:"langkey"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Debug("Failed to decode validate info body: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check the users first name
	if !user.ValidateName(strings.ToLower(data.Firstname)) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeErrorResponse(w, "invalid_first_name")
		return
	}
	// Check the users last name
	if !user.ValidateName(strings.ToLower(data.Lastname)) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeErrorResponse(w, "invalid_last_name")
		return
	}

	counter := 0
	var username string
	for _, r := range data.Firstname {
		if unicode.IsSpace(r) {
			continue
		}
		username += string(unicode.ToLower(r))
	}
	username += "_"
	for _, r := range data.Lastname {
		if unicode.IsSpace(r) {
			continue
		}
		username += string(unicode.ToLower(r))
	}
	username += "_"
	userMgr := user.NewManager(r)

	count, err := userMgr.GetPendingRegistrationsCount()
	if err != nil {
		log.Error("Failed to get pending registerations count: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Debug("count", count)
	if count >= MAX_PENDING_REGISTRATION_COUNT {
		log.Warn("Maximum amount of pending registrations reached")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	orgMgr := organization.NewManager(r)
	exists := true
	for exists {
		counter++
		var err error
		exists, err = userMgr.Exists(username + strconv.Itoa(counter))
		if err != nil {
			log.Error("Failed to verify if username is taken: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !exists {
			exists = orgMgr.Exists(username + strconv.Itoa(counter))
		}
	}
	username = username + strconv.Itoa(counter)

	// Convert the email address to all lowercase
	// Email addresses are limited to printable ASCII characters
	// See https://tools.ietf.org/html/rfc5322#section-3.4.1 for details
	data.Email = strings.ToLower(data.Email)
	valid := user.ValidateEmailAddress(data.Email)
	if !valid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeErrorResponse(w, "invalid_email_format")
		return
	}

	// Check if the email is already known
	valMgr := validationdb.NewManager(r)
	if _, err = valMgr.GetByEmailAddress(data.Email); !db.IsNotFound(err) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeErrorResponse(w, "email_already_used")
		return
	}

	valid = user.ValidatePhoneNumber(data.Phone)
	if !valid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeErrorResponse(w, "invalid_phonenumber")
		return
	}

	// Check if the phone number is already known
	if _, err = valMgr.GetByPhoneNumber(data.Phone); !db.IsNotFound(err) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeErrorResponse(w, "phone_already_used")
		return
	}

	registrationSession, err := service.GetSession(r, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error("Failed to retrieve registration session: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	validatingPhonenumber, _ := registrationSession.Values["phonenumber"].(string)
	validatingEmail, _ := registrationSession.Values["email"].(string)
	validatingUsername, _ := registrationSession.Values["username"].(string)
	validatingPassword, _ := registrationSession.Values["password"].(string)

	phoneChanged := validatingPhonenumber != data.Phone
	emailChanged := validatingEmail != data.Email

	userObj, err := userMgr.GetByName(validatingUsername)
	if err != nil && !db.IsNotFound(err) {
		log.Error("Failed to retrieve user from database: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if db.IsNotFound(err) {
		// If there is no validatingUsername on the request, create a new user
		log.Debug("Creating new user with username ", username)
		userObj = &user.User{
			Username:       username,
			Firstname:      data.Firstname,
			Lastname:       data.Lastname,
			EmailAddresses: []user.EmailAddress{{Label: "main", EmailAddress: data.Email}},
			Phonenumbers:   []user.Phonenumber{{Label: "main", Phonenumber: data.Phone}},
		}

		// give users a day to validate a phone number on their accounts
		duration := time.Duration(time.Hour * 24)
		expiresAt := time.Now()
		expiresAt = expiresAt.Add(duration)
		eat := db.DateTime(expiresAt)
		userObj.Expire = eat
		err = userMgr.Save(userObj)
		if err != nil {
			log.Error("Failed to create new user: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registrationSession.Values["username"] = username
	} else {
		// Update existing user
		// Don't update the username so we don't create a pointer from a password
		// or already confirmed email/phone number to a username that nog longer exists
		userObj.Firstname = data.Firstname
		userObj.Lastname = data.Lastname
		userObj.EmailAddresses = []user.EmailAddress{{Label: "main", EmailAddress: data.Email}}
		userObj.Phonenumbers = []user.Phonenumber{{Label: "main", Phonenumber: data.Phone}}
		err = userMgr.Save(userObj)
		if err != nil {
			log.Error("Failed to update user: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Correct the username
		username = validatingUsername
	}

	if validatingPassword != data.Password || validatingUsername != username {
		log.Debug("Saving user password")
		passwdMgr := password.NewManager(r)
		err = passwdMgr.Save(username, data.Password)
		if err != nil {
			log.Error("Error while saving the users password: ", err)
			if err.Error() != "internal_error" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				writeErrorResponse(w, "invalid_password")
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
		registrationSession.Values["password"] = data.Password
	}

	oldPhoneKey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)
	phoneConfirmed, err := service.phonenumberValidationService.IsConfirmed(r, oldPhoneKey)
	if err != nil && err != validation.ErrInvalidOrExpiredKey {
		log.Error("Failed to check if phone number is already confirmed: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// phone number validation
	// always set the registrationsession value of the phonenumber key. If the phone validation is not triggered,
	// the old value and new value are the same anyway
	registrationSession.Values["phonenumber"] = data.Phone
	if phoneChanged {
		// invalidate old phone number validation
		_ = service.phonenumberValidationService.ExpireValidation(r, oldPhoneKey)

		phonenumber := user.Phonenumber{Phonenumber: data.Phone}
		validationkey, err := service.phonenumberValidationService.RequestValidation(r, username, phonenumber, fmt.Sprintf("https://%s/phonevalidation", r.Host), data.LangKey)
		if err != nil {
			log.Error("Failed to send phonenumber verification in registration flow: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registrationSession.Values["phonenumbervalidationkey"] = validationkey
	}

	// Email validation
	// So the logic here: only send an email if the email address is changed,
	// also only send it if the phone number is confirmed already (defer sending email until this is done)
	// and make sure the phone number didn't change so we don't end up sending 2 validations if
	// the user manages to somehow confirm a wrong phonenumber (magic?)
	// always set the registrationsession value of the email key. If the email validation is not triggered,
	// the old value and new value are the same anyway
	registrationSession.Values["email"] = data.Email
	if emailChanged && phoneConfirmed && !phoneChanged {
		// invalidated old email validation
		oldkey, _ := registrationSession.Values["emailvalidationkey"].(string)
		_ = service.emailaddressValidationService.ExpireValidation(r, oldkey)

		mailvalidationkey, err := service.emailaddressValidationService.RequestValidation(r, username, data.Email, fmt.Sprintf("https://%s/emailvalidation", r.Host), data.LangKey)
		if err != nil {
			log.Error("Failed to send email verification in registration flow: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registrationSession.Values["emailvalidationkey"] = mailvalidationkey
	}

	sessions.Save(r, w)
	// validations created
	w.WriteHeader(http.StatusCreated)
}

func (service *Service) ResendValidationInfo(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		LangKey string `json:"langkey"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Debug("Failed to decode validate info body: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Convert the email to all lowercase
	data.Email = strings.ToLower(data.Email)

	registrationSession, err := service.GetSession(r, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registrationSession.IsNew {
		sessions.Save(r, w)
		log.Debug("Registration session expired")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username, _ := registrationSession.Values["username"].(string)

	// There is no point in resending the validation request if the phone is already
	// verified
	phonevalidationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)
	phoneConfirmed, err := service.phonenumberValidationService.IsConfirmed(r, phonevalidationkey)
	if err != nil {
		log.Error("Failed to check if phone number is already confirmed: ", err)
	}
	if phoneConfirmed {
		log.Debug("Phone is already confirmed, ignoring new phone validation request")
	}
	// Only retrigger the validation if the phone is not confirmed yet
	if !phoneConfirmed && err == nil {

		// Invalidate the previous phone validation request, ignore a possible error
		validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)
		_ = service.phonenumberValidationService.ExpireValidation(r, validationkey)

		phonenumber, _ := registrationSession.Values["phonenumber"].(string)

		validationkey, err = service.phonenumberValidationService.RequestValidation(r, username, user.Phonenumber{Phonenumber: phonenumber}, fmt.Sprintf("https://%s/phonevalidation", r.Host), data.LangKey)
		if err != nil {
			log.Error("ResendPhonenumberConfirmation: Could not get validationkey: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registrationSession.Values["phonenumbervalidationkey"] = validationkey

	}

	// There is no point in resending the validation request if the email is already
	// verified
	emailvalidationkey, _ := registrationSession.Values["emailvalidationkey"].(string)
	emailConfirmed, err := service.emailaddressValidationService.IsConfirmed(r, emailvalidationkey)
	if err != nil && err != validation.ErrInvalidOrExpiredKey {
		log.Error("Failed to check if email address is already confirmed: ", err)
	}
	if emailConfirmed {
		log.Debug("Email is already confirmed, ignoring new email validation request")
	}
	// Only retrigger the validation if the email is not confirmed yet
	// Check if the phone is confirmed, if it is not we can't be on the email page yet
	// So there is no need to resend this validation yet
	if !emailConfirmed && (err == nil || err == validation.ErrInvalidOrExpiredKey) && phoneConfirmed {
		// Invalidate the previous email validation request, ignore a possible error
		_ = service.emailaddressValidationService.ExpireValidation(r, emailvalidationkey)

		email, _ := registrationSession.Values["email"].(string)

		if email != data.Email {
			sessions.Save(r, w)
			log.Info("Attempt to trigger registration flow email (resend) validation with a different email address than the one stored in the session")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		emailvalidationkey, err = service.emailaddressValidationService.RequestValidation(r, username, email, fmt.Sprintf("https://%s/emailvalidation", r.Host), data.LangKey)
		if err != nil {
			log.Error("ResendEmailConfirmation: Could not get validationkey: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registrationSession.Values["emailvalidationkey"] = emailvalidationkey
	}

	sessions.Save(r, w)
	w.WriteHeader(http.StatusOK)
}

func writeErrorResponse(w http.ResponseWriter, err string) {
	response := struct {
		Error string `json:"error"`
	}{
		Error: err,
	}
	json.NewEncoder(w).Encode(&response)
}
