package siteservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/itsyouonline/identityserver/db/persistentlog"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/registration"
	"github.com/itsyouonline/identityserver/db/user"
	validationdb "github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"
	"github.com/itsyouonline/identityserver/tools"
	"github.com/itsyouonline/identityserver/validation"
)

const (
	MAX_PENDING_REGISTRATION_COUNT = 10000
	registrationFileName           = "registration.html"
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
func (service *Service) CheckRegistrationSMSConfirmation(w http.ResponseWriter, r *http.Request) {
	registrationSession, err := service.GetSession(r, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := map[string]bool{}

	if registrationSession.IsNew {
		// TODO: registrationSession is new with SMS, something must be wrong
		log.Info("Registration session timed out while polling for phone confirmation")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	sessionKey, _ := registrationSession.Values["sessionkey"].(string)
	if sessionKey == "" {
		log.Debug("Polling for sms confirmation without session key")
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	registeringUser, err := registration.NewManager(r).GetRegisteringUserBySessionKey(sessionKey)
	if err != nil {
		log.Error("Failed to load registering user object: ", err)
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	confirmed, err := service.phonenumberValidationService.IsConfirmed(r, registeringUser.PhoneValidationKey)
	if err == validation.ErrInvalidOrExpiredKey {
		confirmed = true //This way the form will be submitted, let the form handler deal with redirect to login
		return
	}
	if err != nil {
		log.Error("Failed to check if phone is confirmed in registration flow: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response["confirmed"] = confirmed

	if confirmed {
		persistentlog.NewManager(r).SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "Users phone is confirmed"))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//CheckRegistrationEmailConfirmation is called by the regisration form to check if the email is already confirmed
func (service *Service) CheckRegistrationEmailConfirmation(w http.ResponseWriter, r *http.Request) {
	registrationSession, err := service.GetSession(r, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error("Failed to get registration session: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := map[string]bool{}

	if registrationSession.IsNew {
		// TODO: registrationSession is new, something must be wrong
		log.Warn("Registration is new")
		response["confirmed"] = true //This way the form will be submitted, let the form handler deal with redirect to login
		return
	}

	sessionKey, _ := registrationSession.Values["sessionkey"].(string)
	if sessionKey == "" {
		log.Debug("Polling for sms confirmation without session key")
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	registeringUser, err := registration.NewManager(r).GetRegisteringUserBySessionKey(sessionKey)
	if err != nil {
		log.Error("Failed to load registering user object: ", err)
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	confirmed, err := service.emailaddressValidationService.IsConfirmed(r, registeringUser.EmailValidationKey)
	if err == validation.ErrInvalidOrExpiredKey {
		// TODO
		confirmed = true //This way the form will be submitted, let the form handler deal with redirect to login
		return
	}
	if err != nil {
		log.Error("Failed to check if email is confirmed in registartion flow: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response["confirmed"] = confirmed

	if confirmed {
		persistentlog.NewManager(r).SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "Users email is confirmed"))
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

	if values.Smscode == "" {
		log.Debug("Tried to confirm phone number with empty sms code")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	registrationSession, err := service.GetSession(r, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Debug("Failed to get registration session: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registrationSession.IsNew {
		sessions.Save(r, w)
		log.Debug("Registration session expired")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	sessionKey, _ := registrationSession.Values["sessionkey"].(string)
	if sessionKey == "" {
		log.Debug("Trying to confirm phone number without session key")
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	registeringUser, err := registration.NewManager(r).GetRegisteringUserBySessionKey(sessionKey)
	if err != nil {
		log.Error("Failed to load registering user object: ", err)
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logMgr := persistentlog.NewManager(r)
	if isConfirmed, _ := service.phonenumberValidationService.IsConfirmed(r, registeringUser.PhoneValidationKey); isConfirmed {
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User processed phone number form, but phone is already confirmed"))
		response.Confirmed = true
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = service.phonenumberValidationService.ConfirmRegistrationValidation(r, registeringUser.PhoneValidationKey, values.Smscode)
	if err == validation.ErrInvalidCode {
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User processed phone number form with invalid code"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		response.Error = "invalid_sms_code"
		json.NewEncoder(w).Encode(&response)
		return
	}
	if err == validation.ErrInvalidOrExpiredKey {
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User processed phone number form, but key is expired"))
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&response)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User processed phone number form successfully, phone is confirmed"))
	response.Confirmed = true
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

	sessionKey, _ := registrationSession.Values["sessionkey"].(string)
	if sessionKey == "" {
		log.Debug("Tried to finalize registration but no session key is present")
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	rMgr := registration.NewManager(r)
	registeringUser, err := rMgr.GetRegisteringUserBySessionKey(sessionKey)
	if err != nil {
		log.Error("Failed to get registration object: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logMgr := persistentlog.NewManager(r)

	// check if phone number is validated or sms code is provided to validate phone
	phonevalidationkey := registeringUser.PhoneValidationKey

	if isConfirmed, _ := service.phonenumberValidationService.IsConfirmed(r, phonevalidationkey); !isConfirmed {

		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User tried to register without validated phone"))

		smscode := values.PhonenumberCode
		if smscode == "" {
			log.Debug("no sms code provided and phone not confirmed yet")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = service.phonenumberValidationService.ConfirmRegistrationValidation(r, phonevalidationkey, smscode)
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

		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User confirmed phone number in final registration step"))

	}

	emailvalidationkey := registeringUser.EmailValidationKey
	emailConfirmed, _ := service.emailaddressValidationService.IsConfirmed(r, emailvalidationkey)
	if !emailConfirmed {
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User tried to register without validated email"))
		log.Debug("Email not confirmed yet")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// We know that the phone number and email address are confirmed. So lets create the actual user
	userMgr := user.NewManager(r)

	// First get the username from the firstname and lastname
	username, err := generateUsername(r, registeringUser.Firstname, registeringUser.Lastname)
	if err != nil {
		log.Error("Failed to generate username: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Create the user object
	log.Debug("Creating new user with username ", username)
	userObj := &user.User{
		Username:       username,
		Firstname:      registeringUser.Firstname,
		Lastname:       registeringUser.Lastname,
		EmailAddresses: []user.EmailAddress{{Label: "main", EmailAddress: registeringUser.Email}},
		Phonenumbers:   []user.Phonenumber{{Label: "main", Phonenumber: registeringUser.Phonenumber}},
	}

	err = userMgr.Save(userObj)
	if err != nil {
		log.Error("Failed to create new user: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Also give this user a password
	log.Debug("Saving user password")
	passwdMgr := password.NewManager(r)
	err = passwdMgr.Save(username, values.Password)
	if err != nil {
		log.Error("Error while saving the users password: ", err)
		if err.Error() != "internal_error" {
			writeErrorResponse(w, "invalid_password", http.StatusUnprocessableEntity)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Add correct validated phone number and validated email address
	valMgr := validationdb.NewManager(r)
	p := valMgr.NewValidatedPhonenumber(username, registeringUser.Phonenumber)
	if err = valMgr.SaveValidatedPhonenumber(p); err != nil {
		log.Error("Failed to add validated phone number of new user: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	e := valMgr.NewValidatedEmailAddress(username, registeringUser.Email)
	if err = valMgr.SaveValidatedEmailAddress(e); err != nil {
		log.Error("Failed to add validated email address of new user: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "Registration finished, user created"))
	log.Debug("Finished saving new user information")

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

// ValidateInfo starts validation for a temporary username
func (service *Service) ValidateInfo(w http.ResponseWriter, r *http.Request) {
	registrationSession, err := service.GetSession(r, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Error("Failed to retrieve registration session: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rMgr := registration.NewManager(r)

	logMgr := persistentlog.NewManager(r)

	var registeringUser *registration.InProgressRegistration
	sessionKey, existingSession := registrationSession.Values["sessionkey"].(string)
	if existingSession {
		registeringUser, err = rMgr.GetRegisteringUserBySessionKey(sessionKey)
		if err != nil {
			log.Error("Failed to get registration object: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User updated his info"))
	} else {
		sessionKey, err = tools.GenerateRandomString()
		if err != nil {
			log.Error("Failed to generate session key: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registeringUser = registration.New(sessionKey)
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, fmt.Sprintf("User started registration flow. State: %s", r.URL.Query().Get("state"))))
	}

	registrationSession.Values["sessionkey"] = sessionKey

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
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User first name invalid"))
		writeErrorResponse(w, "invalid_first_name", http.StatusUnprocessableEntity)
		return
	}
	registeringUser.Firstname = data.Firstname

	// Check the users last name
	if !user.ValidateName(strings.ToLower(data.Lastname)) {
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User last name invalid"))
		writeErrorResponse(w, "invalid_last_name", http.StatusUnprocessableEntity)
		return
	}
	registeringUser.Lastname = data.Lastname

	// Convert the email address to all lowercase
	// Email addresses are limited to printable ASCII characters
	// See https://tools.ietf.org/html/rfc5322#section-3.4.1 for details
	data.Email = strings.ToLower(data.Email)
	if !user.ValidateEmailAddress(data.Email) {
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User email invalid"))
		writeErrorResponse(w, "invalid_email_format", http.StatusUnprocessableEntity)
		return
	}

	// Check if the email is already known
	valMgr := validationdb.NewManager(r)
	if _, err = valMgr.GetByEmailAddress(data.Email); !db.IsNotFound(err) {
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User email already confirmed by someone else"))
		writeErrorResponse(w, "email_already_used", http.StatusUnprocessableEntity)
		return
	}
	newEmail := data.Email != registeringUser.Email
	registeringUser.Email = data.Email

	if !user.ValidatePhoneNumber(data.Phone) {
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User phone number invalid"))
		writeErrorResponse(w, "invalid_phonenumber", http.StatusUnprocessableEntity)
		return
	}

	// Check if the phone number is already known
	if _, err = valMgr.GetByPhoneNumber(data.Phone); !db.IsNotFound(err) {
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User phone number already confirmed by someone else"))
		writeErrorResponse(w, "phone_already_used", http.StatusUnprocessableEntity)
		return
	}
	newPhone := data.Phone != registeringUser.Phonenumber
	registeringUser.Phonenumber = data.Phone

	// Check the password
	if err = password.Check(data.Password); err != nil {
		if err != password.ErrInvalidPassword {
			log.Error("Failed to verify password: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User password invalid"))
		log.Debug("User password is invalid")
		writeErrorResponse(w, "invalid_password", http.StatusUnprocessableEntity)
		return
	}
	// Storing password in plaintext is probably not such a great idea
	// Therefore, just do a check now, and do the password thing when we actually
	// create the user

	phoneConfirmed, err := service.phonenumberValidationService.IsConfirmed(r, registeringUser.PhoneValidationKey)
	if err != nil && err != validation.ErrInvalidOrExpiredKey {
		log.Error("Failed to check if phone number is already confirmed: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// phone number validation
	if newPhone && !phoneConfirmed {
		// invalidate old phone number validation
		_ = service.phonenumberValidationService.ExpireValidation(r, registeringUser.PhoneValidationKey)

		phonenumber := user.Phonenumber{Phonenumber: data.Phone}
		validationkey, err := service.phonenumberValidationService.RequestValidation(r, registeringUser.SessionKey, phonenumber, fmt.Sprintf("https://%s/phoneregistrationvalidation", r.Host), data.LangKey)
		if err != nil {
			log.Error("Failed to send phonenumber verification in registration flow: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registeringUser.PhoneValidationKey = validationkey

		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User started phone validation for new phone"))
	}

	// Email validation
	// So the logic here: only send an email if the email address is changed,
	// also only send it if the phone number is confirmed already (defer sending email until this is done)
	if newEmail && phoneConfirmed {
		// invalidated old email validation
		_ = service.emailaddressValidationService.ExpireValidation(r, registeringUser.EmailValidationKey)

		mailvalidationkey, err := service.emailaddressValidationService.RequestValidation(r, registeringUser.SessionKey, data.Email, fmt.Sprintf("https://%s/emailregistrationvalidation", r.Host), data.LangKey)
		if err != nil {
			log.Error("Failed to send email verification in registration flow: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registeringUser.EmailValidationKey = mailvalidationkey

		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User started email validation for new email"))
	}

	// Save the info we have so far about the user registering
	if err = rMgr.UpsertRegisteringUser(registeringUser); err != nil {
		log.Error("Failed to save registering user to database: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "User info processed, moving to validations"))

	sessions.Save(r, w)
	// validations created
	w.WriteHeader(http.StatusCreated)

}

// ResendValidationInfo resends validation info for either the phone number or email address
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

	sessionKey, _ := registrationSession.Values["sessionkey"].(string)
	if sessionKey == "" {
		log.Debug("Requested to resend registration info but no session key is present")
		sessions.Save(r, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	logMgr := persistentlog.NewManager(r)

	rMgr := registration.NewManager(r)
	registeringUser, err := rMgr.GetRegisteringUserBySessionKey(sessionKey)
	if err != nil {
		log.Error("Failed to get registration object: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// There is no point in resending the validation request if the phone is already
	// verified
	phonevalidationkey := registeringUser.PhoneValidationKey
	phoneConfirmed, err := service.phonenumberValidationService.IsConfirmed(r, phonevalidationkey)
	if err != nil {
		log.Error("Failed to check if phone number is already confirmed: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !phoneConfirmed {
		// Invalidate the previous phone validation request, ignore a possible error
		_ = service.phonenumberValidationService.ExpireValidation(r, registeringUser.PhoneValidationKey)

		phonenumber := registeringUser.Phonenumber

		if phonenumber != data.Phone {
			sessions.Save(r, w)
			log.Info("Attempt to trigger registration flow phone (resend) validation with a different phone number than the one stored in the session")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Debug("Sending new phone number confirmation")

		validationkey, err := service.phonenumberValidationService.RequestValidation(r, registeringUser.SessionKey, user.Phonenumber{Phonenumber: phonenumber}, fmt.Sprintf("https://%s/phoneregistrationvalidation", r.Host), data.LangKey)
		if err != nil {
			log.Error("ResendPhonenumberConfirmation: Could not get validationkey: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registeringUser.PhoneValidationKey = validationkey

		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "Resending phone validation"))
	}

	// There is no point in resending the validation request if the email is already
	// verified
	emailvalidationkey := registeringUser.EmailValidationKey
	emailConfirmed, err := service.emailaddressValidationService.IsConfirmed(r, emailvalidationkey)
	if err != nil && err != validation.ErrInvalidOrExpiredKey {
		log.Error("Failed to check if email address is already confirmed: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if phoneConfirmed && !emailConfirmed {
		// Invalidate the previous email validation request, ignore a possible error
		_ = service.emailaddressValidationService.ExpireValidation(r, registeringUser.EmailValidationKey)

		email := registeringUser.Email

		if email != data.Email {
			sessions.Save(r, w)
			log.Info("Attempt to trigger registration flow email (resend) validation with a different email address than the one stored in the session")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Debug("Sending new email validation")
		emailvalidationkey, err := service.emailaddressValidationService.RequestValidation(r, sessionKey, email, fmt.Sprintf("https://%s/emailregistrationvalidation", r.Host), data.LangKey)
		if err != nil {
			log.Error("ResendEmailConfirmation: Could not get validationkey: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registeringUser.EmailValidationKey = emailvalidationkey

		logMgr.SaveLog(persistentlog.New(sessionKey, persistentlog.RegistrationFlow, "Resending emal validation"))
	}

	// Save the info we have so far about the user registering
	if err = rMgr.UpsertRegisteringUser(registeringUser); err != nil {
		log.Error("Failed to save registering user to database: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sessions.Save(r, w)
	w.WriteHeader(http.StatusOK)
}

// generateUsername generates a new username
func generateUsername(r *http.Request, firstname, lastname string) (string, error) {
	counter := 0
	var username string
	for _, r := range firstname {
		if unicode.IsSpace(r) {
			continue
		}
		username += string(unicode.ToLower(r))
	}
	username += "_"
	for _, r := range lastname {
		if unicode.IsSpace(r) {
			continue
		}
		username += string(unicode.ToLower(r))
	}
	username += "_"
	userMgr := user.NewManager(r)

	count, err := userMgr.GetPendingRegistrationsCount()
	if err != nil {
		return "", err
	}
	log.Debug("count", count)
	if count >= MAX_PENDING_REGISTRATION_COUNT {
		return "", errors.New("Max amount of pending registrations reached")
	}

	orgMgr := organization.NewManager(r)
	exists := true
	for exists {
		counter++
		var err error
		exists, err = userMgr.Exists(username + strconv.Itoa(counter))
		if err != nil {
			return "", err
		}
		if !exists {
			exists = orgMgr.Exists(username + strconv.Itoa(counter))
		}
	}
	username = username + strconv.Itoa(counter)
	return username, nil
}

func writeErrorResponse(w http.ResponseWriter, err string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := struct {
		Error string `json:"error"`
	}{
		Error: err,
	}
	json.NewEncoder(w).Encode(&response)
}
