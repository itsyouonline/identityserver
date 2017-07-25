package siteservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/credentials/totp"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/user"
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
		// todo: registrationSession is new with SMS, something must be wrong
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

//ShowRegistrationForm shows the user registration page
func (service *Service) ShowRegistrationForm(w http.ResponseWriter, request *http.Request) {
	service.renderRegistrationFrom(w, request)
}

//ProcessPhonenumberConfirmationForm processes the Phone number confirmation form
func (service *Service) ProcessPhonenumberConfirmationForm(w http.ResponseWriter, request *http.Request) {
	values := struct {
		Smscode string `json:"smscode"`
	}{}

	response := struct {
		RedirectUrL string `json:"redirecturl"`
		Error       string `json:"error"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ProcessPhonenumberConfirmation request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	registrationSession, err := service.GetSession(request, SessionForRegistration, "registrationdetails")
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registrationSession.IsNew {
		sessions.Save(request, w)
		response.RedirectUrL = fmt.Sprintf("https://%s/register/#/smsconfirmation", request.Host)
		json.NewEncoder(w).Encode(&response)
		return
	}

	username, _ := registrationSession.Values["username"].(string)
	validationkey, _ := registrationSession.Values["phonenumbervalidationkey"].(string)

	if isConfirmed, _ := service.phonenumberValidationService.IsConfirmed(request, validationkey); isConfirmed {
		userMgr := user.NewManager(request)
		userMgr.RemoveExpireDate(username)
		service.loginUser(w, request, username)
		return
	}

	smscode := values.Smscode
	if err != nil || smscode == "" {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = service.phonenumberValidationService.ConfirmValidation(request, validationkey, smscode)
	if err == validation.ErrInvalidCode {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response.Error = "invalid_sms_code"
		json.NewEncoder(w).Encode(&response)
		return
	}
	if err == validation.ErrInvalidOrExpiredKey {
		sessions.Save(request, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&response)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	userMgr := user.NewManager(request)
	userMgr.RemoveExpireDate(username)
	service.loginUser(w, request, username)
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
func (service *Service) ProcessRegistrationForm(w http.ResponseWriter, request *http.Request) {
	response := struct {
		Redirecturl string `json:"redirecturl"`
		Error       string `json:"error"`
	}{}
	values := struct {
		TwoFAMethod    string `json:"twofamethod"`
		Login          string `json:"login"`
		Email          string `json:"email"`
		Phonenumber    string `json:"phonenumber"`
		TotpCode       string `json:"totpcode"`
		Password       string `json:"password"`
		RedirectParams string `json:"redirectparams"`
		LangKey        string `json:"langkey"`
	}{}
	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the registration request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	twoFAMethod := values.TwoFAMethod
	if twoFAMethod != "sms" && twoFAMethod != "totp" {
		log.Info("Invalid 2fa method during registration: ", twoFAMethod)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	totpsession, err := service.GetSession(request, SessionForRegistration, "totp")
	if err != nil {
		log.Error("ERROR while getting the totp registration session", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if totpsession.IsNew {
		log.Debug("New registration session while processing the registration form")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	totpsecret, ok := totpsession.Values["secret"].(string)
	if !ok {
		log.Error("Unable to convert the stored session totp secret to a string")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	valid := user.ValidateUsername(values.Login)
	var phonenumber user.Phonenumber
	if !valid {
		response.Error = "invalid_username_format"
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response)
		return
	}

	valid = user.ValidateEmailAddress(values.Email)
	if !valid {
		response.Error = "invalid_email_format"
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response)
		return
	}

	newuser := &user.User{
		Username:       values.Login,
		EmailAddresses: []user.EmailAddress{{Label: "main", EmailAddress: values.Email}},
	}
	//validate the username is not taken yet
	userMgr := user.NewManager(request)
	orgMgr := organization.NewManager(request)

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

	//we now just depend on mongo unique index to avoid duplicates when concurrent requests are made
	userExists, err := userMgr.Exists(newuser.Username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if userExists {
		log.Debug("USER ", newuser.Username, " already registered")
		w.WriteHeader(http.StatusUnprocessableEntity)
		response.Error = "user_exists"
		json.NewEncoder(w).Encode(&response)
		return
	} else if orgMgr.Exists(newuser.Username) {
		log.Debugf("Cannot create user: organization with globalid %s already exists", newuser.Username)
		w.WriteHeader(http.StatusUnprocessableEntity)
		response.Error = "organization_exists"
		json.NewEncoder(w).Encode(&response)
		return
	}

	if twoFAMethod == "sms" {
		phonenumber = user.Phonenumber{Label: "main", Phonenumber: values.Phonenumber}
		if !phonenumber.Validate() {
			log.Debug("Invalid phone number")
			w.WriteHeader(http.StatusUnprocessableEntity)
			response.Error = "invalid_phonenumber"
			json.NewEncoder(w).Encode(&response)
			return
		}
		newuser.Phonenumbers = []user.Phonenumber{phonenumber}
		// Remove account after 3 days if it still doesn't have a verified phone by then
		duration := time.Duration(time.Hour * 24 * 3)
		expiresAt := time.Now()
		expiresAt = expiresAt.Add(duration)
		newuser.Expire = db.DateTime(expiresAt)
	} else {
		token := totp.TokenFromSecret(totpsecret)
		if !token.Validate(values.TotpCode) {
			log.Debug("Invalid totp code")
			w.WriteHeader(http.StatusUnprocessableEntity)
			response.Error = "invalid_totpcode"
			json.NewEncoder(w).Encode(&response)
			return
		}
	}

	userMgr.Save(newuser)
	passwdMgr := password.NewManager(request)
	err = passwdMgr.Save(newuser.Username, values.Password)
	if err != nil {
		log.Error(err)
		if err.Error() != "internal_error" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			response.Error = "invalid_password"
			json.NewEncoder(w).Encode(&response)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// send an email to ask for email address validation now the user is saved
	_, err = service.emailaddressValidationService.RequestValidation(request, newuser.Username, values.Email, fmt.Sprintf("https://%s/emailvalidation", request.Host), values.LangKey)
	if err != nil {
		// Failure to send the email is an error, but not critical so don't abort the flow
		log.Error("Failed to send email verification in registration flow: ", err)
	}

	if twoFAMethod == "sms" {
		validationkey, err := service.phonenumberValidationService.RequestValidation(request, newuser.Username, phonenumber, fmt.Sprintf("https://%s/phonevalidation", request.Host), values.LangKey)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		registrationSession, err := service.GetSession(request, SessionForRegistration, "registrationdetails")
		registrationSession.Values["username"] = newuser.Username
		registrationSession.Values["phonenumbervalidationkey"] = validationkey
		registrationSession.Values["redirectparams"] = values.RedirectParams

		sessions.Save(request, w)
		response.Redirecturl = fmt.Sprintf("https://%s/register?%s#/smsconfirmation", request.Host, values.RedirectParams)
		json.NewEncoder(w).Encode(&response)
		return
	}

	totpMgr := totp.NewManager(request)
	totpMgr.Save(newuser.Username, totpsecret)
	log.Debugf("Registered %s", newuser.Username)
	service.loginUser(w, request, newuser.Username)
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
