package siteservice

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"
	"gopkg.in/mgo.v2"

	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/credentials/totp"
	"github.com/itsyouonline/identityserver/db/user"
	organizationdb "github.com/itsyouonline/identityserver/db/organization"
	validationdb "github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/identityservice/organization"
	"github.com/itsyouonline/identityserver/tools"
	"github.com/itsyouonline/identityserver/validation"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoLoginCollectionName = "loginsessions"
)

//initLoginModels initialize models in mongo
func (service *Service) initLoginModels() {
	index := mgo.Index{
		Key:      []string{"sessionkey"},
		Unique:   true,
		DropDups: false,
	}

	db.EnsureIndex(mongoLoginCollectionName, index)

	automaticExpiration := mgo.Index{
		Key:         []string{"createdat"},
		ExpireAfter: time.Second * 60 * 10,
		Background:  true,
	}
	db.EnsureIndex(mongoLoginCollectionName, automaticExpiration)

}

type loginSessionInformation struct {
	SessionKey string
	SMSCode    string
	Confirmed  bool
	CreatedAt  time.Time
}

func newLoginSessionInformation() (sessionInformation *loginSessionInformation, err error) {
	sessionInformation = &loginSessionInformation{CreatedAt: time.Now()}
	sessionInformation.SessionKey, err = tools.GenerateRandomString()
	if err != nil {
		return
	}
	numbercode, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return
	}
	sessionInformation.SMSCode = fmt.Sprintf("%06d", numbercode)
	return
}

const loginFileName = "login.html"

//ShowLoginForm shows the user login page on the initial request
func (service *Service) ShowLoginForm(w http.ResponseWriter, request *http.Request) {
	htmlData, err := html.Asset(loginFileName)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loginSession.Values["auth_client_id"] = request.URL.Query().Get("client_id")
	sessions.Save(request, w)
	w.Write(htmlData)

}

//ProcessLoginForm logs a user in if the credentials are valid
func (service *Service) ProcessLoginForm(w http.ResponseWriter, request *http.Request) {
	//TODO: validate csrf token
	//TODO: limit the number of failed/concurrent requests

	err := request.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing registration form")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	values := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}

	if err = json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the login request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	login := strings.ToLower(values.Login)

	u, err := organization.SearchUser(request, login)
	if err == mgo.ErrNotFound {
		w.WriteHeader(422)
		return
	} else if err != nil {
		log.Error("Failed to search for user: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	userexists := err != mgo.ErrNotFound

	var validpassword bool
	passwdMgr := password.NewManager(request)
	if validpassword, err = passwdMgr.Validate(u.Username, values.Password); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	queryValues := request.URL.Query()
	client := queryValues.Get("client_id")
	// Remove last 2FA entry if an invalid password is entered
	validcredentials := userexists && validpassword
	if !validcredentials {
		if client != "" {
			l2faMgr := organizationdb.NewLast2FAManager(request)
			if l2faMgr.Exists(client, u.Username) {
				l2faMgr.RemoveLast2FA(client, u.Username)
			}
		}
		w.WriteHeader(422)
		return
	}
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loginSession.Values["username"] = u.Username
	//check if 2fa validity has passed
	if client != "" {
		l2faMgr := organizationdb.NewLast2FAManager(request)
		if l2faMgr.Exists(client, u.Username) {
			timestamp, err := l2faMgr.GetLast2FA(client, u.Username)
			if err != nil {
				log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			mgr := organizationdb.NewManager(request)
			seconds, err := mgr.GetValidity(client)
			if err != nil {
				log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			timeconverted := time.Time(timestamp)
			if timeconverted.Add(time.Second * time.Duration(seconds)).After(time.Now()) {
				service.loginUser(w, request, u.Username)
				return
			}
		}
	}

	sessions.Save(request, w)
	w.WriteHeader(http.StatusNoContent)
}

// GetTwoFactorAuthenticationMethods returns the possible two factor authentication methods the user can use to login with.
func (service *Service) GetTwoFactorAuthenticationMethods(w http.ResponseWriter, request *http.Request) {
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username, ok := loginSession.Values["username"].(string)
	if username == "" || !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	userMgr := user.NewManager(request)
	userFromDB, err := userMgr.GetByName(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := struct {
		Totp bool              `json:"totp"`
		Sms  map[string]string `json:"sms"`
	}{Sms: make(map[string]string)}
	totpMgr := totp.NewManager(request)
	response.Totp, err = totpMgr.HasTOTP(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	valMgr := validationdb.NewManager(request)
	verifiedPhones, err := valMgr.GetByUsernameValidatedPhonenumbers(username)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for _, validatedPhoneNumber := range verifiedPhones {
		for _, number := range userFromDB.Phonenumbers {
			if number.Phonenumber == string(validatedPhoneNumber.Phonenumber) {
				response.Sms[number.Label] = string(validatedPhoneNumber.Phonenumber)
			}
		}
	}
	json.NewEncoder(w).Encode(response)
	return
}

//getUserLoggingIn returns an user trying to log in, or an empty string if there is none
func (service *Service) getUserLoggingIn(request *http.Request) (username string, err error) {
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error(err)
		return
	}
	savedusername := loginSession.Values["username"]
	if savedusername != nil {
		username, _ = savedusername.(string)
	}
	return
}

//getSessionKey returns an the login session key, or an empty string if there is none
func (service *Service) getSessionKey(request *http.Request) (sessionKey string, err error) {
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error(err)
		return
	}
	savedSessionKey := loginSession.Values["sessionkey"]
	if savedSessionKey != nil {
		sessionKey, _ = savedSessionKey.(string)
	}
	return
}

//GetSmsCode returns an sms code for a specified phone label
func (service *Service) GetSmsCode(w http.ResponseWriter, request *http.Request) {
	phoneLabel := mux.Vars(request)["phoneLabel"]
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error("Error getting login session", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sessionInfo, err := newLoginSessionInformation()
	if err != nil {
		log.Error("Error creating login session information", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username, ok := loginSession.Values["username"].(string)
	if username == "" || !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	userMgr := user.NewManager(request)
	userFromDB, err := userMgr.GetByName(username)
	if err != nil {
		log.Error("Error getting user", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	phoneNumber, err := userFromDB.GetPhonenumberByLabel(phoneLabel)
	if err != nil {
		log.Debug(userFromDB.Phonenumbers)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	loginSession.Values["sessionkey"] = sessionInfo.SessionKey
	authClientId := loginSession.Values["auth_client_id"]
	authenticatingOrganization := ""
	if authClientId != nil {
		authenticatingOrganization = authClientId.(string)
	}
	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoLoginCollectionName)
	mgoCollection.Insert(sessionInfo)
	organizationText := ""
	if authenticatingOrganization != "" {
		split := strings.Split(authenticatingOrganization, ".")
		organizationText = fmt.Sprintf("to authorize the organization %s, ", split[len(split)-1])
	}
	smsmessage := fmt.Sprintf("To continue signing in at itsyou.online %senter the code %s in the form or use this link: https://%s/sc?c=%s&k=%s",
		organizationText, sessionInfo.SMSCode, request.Host, sessionInfo.SMSCode, url.QueryEscape(sessionInfo.SessionKey))
	sessions.Save(request, w)
	go service.smsService.Send(phoneNumber.Phonenumber, smsmessage)
	w.WriteHeader(http.StatusNoContent)
}

//ProcessTOTPConfirmation checks the totp 2 factor authentication code
func (service *Service) ProcessTOTPConfirmation(w http.ResponseWriter, request *http.Request) {
	username, err := service.getUserLoggingIn(request)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if username == "" {
		sessions.Save(request, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	values := struct {
		Totpcode string `json:"totpcode"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the totp confirmation request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var validtotpcode bool
	totpMgr := totp.NewManager(request)
	if validtotpcode, err = totpMgr.Validate(username, values.Totpcode); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !validtotpcode { //TODO: limit to 3 failed attempts
		w.WriteHeader(422)
		return
	}

	//add last 2fa date if logging in with oauth2
	service.storeLast2FALogin(request, username)

	service.loginUser(w, request, username)
}

func (service *Service) getLoginSessionInformation(request *http.Request, sessionKey string) (sessionInfo *loginSessionInformation, err error) {

	if sessionKey == "" {
		sessionKey, err = service.getSessionKey(request)
		if err != nil || sessionKey == "" {
			return
		}
	}

	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoLoginCollectionName)
	sessionInfo = &loginSessionInformation{}
	err = mgoCollection.Find(bson.M{"sessionkey": sessionKey}).One(sessionInfo)
	if err == mgo.ErrNotFound {
		sessionInfo = nil
		err = nil
	}
	return
}

//MobileSMSConfirmation is the page that is linked to in the SMS and is thus accessed on the mobile phone
func (service *Service) MobileSMSConfirmation(w http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing mobile smsconfirmation form", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	values := request.Form
	sessionKey := values.Get("k")
	smscode := values.Get("c")

	var validsmscode bool
	sessionInfo, err := service.getLoginSessionInformation(request, sessionKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if sessionInfo == nil {
		service.renderSMSConfirmationPage(w, request, "Invalid or expired link")
		return
	}

	validsmscode = (smscode == sessionInfo.SMSCode)

	if !validsmscode { //TODO: limit to 3 failed attempts
		service.renderSMSConfirmationPage(w, request, "Invalid or expired link")
		return
	}
	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoLoginCollectionName)

	_, err = mgoCollection.UpdateAll(bson.M{"sessionkey": sessionKey}, bson.M{"$set": bson.M{"confirmed": true}})
	if err != nil {
		log.Error("Error while confirming sms 2fa - ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	service.renderSMSConfirmationPage(w, request, "You will be logged in within a few seconds")
}

//Check2FASMSConfirmation is called by the sms code form to check if the sms is already confirmed on the mobile phone
func (service *Service) Check2FASMSConfirmation(w http.ResponseWriter, request *http.Request) {

	sessionInfo, err := service.getLoginSessionInformation(request, "")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := map[string]bool{}
	if sessionInfo == nil {
		response["confirmed"] = false
	} else {
		response["confirmed"] = sessionInfo.Confirmed
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

//Process2FASMSConfirmation checks the totp 2 factor authentication code
func (service *Service) Process2FASMSConfirmation(w http.ResponseWriter, request *http.Request) {
	username, err := service.getUserLoggingIn(request)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if username == "" {
		sessions.Save(request, w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	values := struct {
		Smscode string `json:"smscode"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the totp confirmation request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sessionInfo, err := service.getLoginSessionInformation(request, "")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if sessionInfo == nil {
		loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
		if err != nil {
			if err == mgo.ErrNotFound {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		validationkey, _ := loginSession.Values["phonenumbervalidationkey"].(string)
		err = service.phonenumberValidationService.ConfirmValidation(request, validationkey, values.Smscode)
		if err == validation.ErrInvalidCode {
			// TODO: limit to 3 failed attempts
			w.WriteHeader(422)
			log.Debug("invalid code")
			return
		}
	} else if !sessionInfo.Confirmed {
		//Already confirmed on the phone
		validsmscode := (values.Smscode == sessionInfo.SMSCode)

		if !validsmscode {
			// TODO: limit to 3 failed attempts
			w.WriteHeader(422)
			log.Debugf("Expected code %s, got %s", sessionInfo.SMSCode, values.Smscode)
			return
		}
	}
	userMgr := user.NewManager(request)
	userMgr.RemoveExpireDate(username)

	//add last 2fa date if logging in with oauth2
	service.storeLast2FALogin(request, username)

	service.loginUser(w, request, username)
}

func (service *Service) storeLast2FALogin(request *http.Request, username string) {
	//add last 2fa date if logging in with oauth2
	queryValues := request.URL.Query()
	client := queryValues.Get("client_id")
	if client != "" {
		l2faMgr := organizationdb.NewLast2FAManager(request)
		err := l2faMgr.SetLast2FA(client, username)
		if err != nil {
			log.Debug("Error while setting the last 2FA login ", err)
		}
	}
}

func (service *Service) loginUser(w http.ResponseWriter, request *http.Request, username string) {
	if err := service.SetLoggedInUser(w, request, username); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sessions.Save(request, w)

	log.Debugf("Successfull login by '%s'", username)

	redirectURL := "/"
	queryValues := request.URL.Query()
	endpoint := queryValues.Get("endpoint")
	if endpoint != "" {
		queryValues.Del("endpoint")
		redirectURL = endpoint + "?" + queryValues.Encode()
	} else {
		registrationSession, _ := service.GetSession(request, SessionForRegistration, "registrationdetails")
		if !registrationSession.IsNew && registrationSession.Values["redirectparams"] != nil {
			splitted := strings.Split(registrationSession.Values["redirectparams"].(string), "&")
			if len(splitted) > 3 {
				for _, part := range splitted {
					kv := strings.Split(part, "=")
					if len(kv) == 2 {
						key, _ := url.QueryUnescape(kv[0])
						value, _ := url.QueryUnescape(kv[1])
						queryValues.Set(key, value)
					}
				}
				endpoint, _ = url.QueryUnescape(queryValues.Get("endpoint"))
				queryValues.Del("endpoint")
				redirectURL = endpoint + "?" + queryValues.Encode()
			}
		}
	}

	sessions.Save(request, w)
	response := struct {
		Redirecturl string `json:"redirecturl"`
	}{}
	response.Redirecturl = redirectURL
	log.Debug("Redirecting to:", redirectURL)
	json.NewEncoder(w).Encode(response)
}

//ForgotPassword handler for POST /login/forgotpassword
func (service *Service) ForgotPassword(w http.ResponseWriter, request *http.Request) {
	// login can be username or email
	values := struct {
		Login string `json:"login"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ForgotPassword request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	userMgr := user.NewManager(request)
	valMgr := validationdb.NewManager(request)
	validatedemail, err := valMgr.GetByEmailAddressValidatedEmailAddress(values.Login)
	if err != nil && err != mgo.ErrNotFound {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	var username string
	var emails []string
	if err != mgo.ErrNotFound {
		username = validatedemail.Username
		emails = []string{validatedemail.EmailAddress}
	} else {
		user, err := userMgr.GetByName(values.Login)
		if err != nil && err != mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		username = user.Username
		validatedemails, err := valMgr.GetByUsernameValidatedEmailAddress(username)
		if validatedemails == nil || len(validatedemails) == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err != nil {
			log.Error("Failed to get validated emails address - ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		emails = make([]string, len(validatedemails))
		for idx, validatedemail := range validatedemails {
			emails[idx] = validatedemail.EmailAddress
		}

	}
	_, err = service.emailaddressValidationService.RequestPasswordReset(request, username, emails)
	if err != nil {
		log.Error("Failed to request password reset - ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
	return
}

//ResetPassword handler for POST /login/resetpassword
func (service *Service) ResetPassword(w http.ResponseWriter, request *http.Request) {
	values := struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ResetPassword request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	pwdMngr := password.NewManager(request)
	token, err := pwdMngr.FindResetToken(values.Token)
	if err != nil {
		log.Debug("Failed to find password reset token - ", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	err = pwdMngr.Save(token.Username, values.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err = pwdMngr.DeleteResetToken(values.Token); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusNoContent)
	return
}

//LoginResendPhonenumberConfirmation resend the phone number confirmation after logging in to a possibly new phone number
func (service *Service) LoginResendPhonenumberConfirmation(w http.ResponseWriter, request *http.Request) {
	values := struct {
		PhoneNumber string `json:"phonenumber"`
	}{}

	response := struct {
		Error string `json:"error"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ResendPhonenumberConfirmation request: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if loginSession.IsNew {
		sessions.Save(request, w)
		log.Debug("Login session expired")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username, _ := loginSession.Values["username"].(string)

	//Invalidate the previous validation request, ignore a possible error
	validationkey, _ := loginSession.Values["phonenumbervalidationkey"].(string)
	_ = service.phonenumberValidationService.ExpireValidation(request, validationkey)

	phonenumber := user.Phonenumber{Label: "main", Phonenumber: values.PhoneNumber}
	if !phonenumber.IsValid() {
		log.Debug("Invalid phone number")
		w.WriteHeader(422)
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
	validationkey, err = service.phonenumberValidationService.RequestValidation(request, username, phonenumber, fmt.Sprintf("https://%s/phonevalidation", request.Host))
	if err != nil {
		log.Error("ResendPhonenumberConfirmation: Could not get validationkey: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loginSession.Values["phonenumbervalidationkey"] = validationkey
	sessions.Save(request, w)
	w.WriteHeader(http.StatusNoContent)
}
