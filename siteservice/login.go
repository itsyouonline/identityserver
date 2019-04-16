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
	"github.com/itsyouonline/identityserver/oauthservice"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"
	"gopkg.in/mgo.v2"

	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/credentials/oauth2"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/credentials/totp"
	organizationdb "github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/user"
	validationdb "github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/identityservice/invitations"
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

		// Check if we have a valid authorization
		requestedScopes := oauth2.SplitScopeString(request.Form.Get("scope"))
		possibleScopes, err := service.identityService.FilterPossibleScopes(request, u.Username, requestedScopes, true)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		validAuthorization, err := service.verifyExistingAuthorization(request, u.Username, client, possibleScopes)
		if err != nil {
			log.Error("Failed to check if authorization is valid: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Only attempt to bypass 2fa if we have a valid authorization
		if validAuthorization {
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
					log.Debug("Try to build protected session")
					service.loginOauthUser(w, request, u.Username)
					return
				}
			}
		}
	}

	sessions.Save(request, w)
	w.WriteHeader(http.StatusNoContent)
}

func (service *Service) verifyExistingAuthorization(request *http.Request, username string, clientID string, possibleScopes []string) (bool, error) {
	authorizedScopes, err := service.identityService.FilterAuthorizedScopes(request, username, clientID, possibleScopes)
	if err != nil {
		log.Error(err)
		return false, err
	}

	var validAuthorization bool

	if authorizedScopes != nil {

		// Check if all authorizatons are given
		validAuthorization = oauthservice.IsAuthorizationValid(possibleScopes, authorizedScopes)
		// Check if the user still has the given authorizations
		authorization, err := user.NewManager(request).GetAuthorization(username, clientID)
		if err != nil {
			log.Error("Failed to load authorization: ", err)
			return false, err
		}
		if validAuthorization {
			validAuthorization, err = oauthservice.UserHasAuthorizedScopes(request, authorization)
			if err != nil {
				log.Error("Failed to check if authorizated labels are still present: ", err)
				return false, err
			}
		}

		//Check if we are redirected from the authorize page, it might be that not all authorizations were given,
		// authorize the login but only with the authorized scopes
		referrer := request.Header.Get("Referer")
		if referrer != "" && !validAuthorization { //If we already have a valid authorization, no need to check if we come from the authorize page
			if referrerURL, e := url.Parse(referrer); e == nil {
				validAuthorization = referrerURL.Host == request.Host && referrerURL.Path == "/authorize"
			} else {
				log.Debug("Error parsing referrer: ", e)
			}
		}
	}
	return validAuthorization, err
}

// GetTwoFactorAuthenticationSettings returns the Two FA settings for current user.
func (service *Service) GetTwoFactorAuthenticationSettings(w http.ResponseWriter, request *http.Request) {
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

	json.NewEncoder(w).Encode(userFromDB.TwoFA)
	return
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

	values := struct {
		LangKey string `json:"langkey"`
	}{}
	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the GetSmsCode langkey request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

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

	link := fmt.Sprintf("https://%s/sc?c=%s&k=%s&l=%s", request.Host, sessionInfo.SMSCode, url.QueryEscape(sessionInfo.SessionKey), values.LangKey)
	translationValues := make(tools.TranslationValues)
	if authenticatingOrganization != "" {
		split := strings.Split(authenticatingOrganization, ".")
		translationValues["authorizeorganizationsms"] = struct {
			Organization string
			Code         string
			Link         string
		}{
			Organization: split[len(split)-1],
			Code:         sessionInfo.SMSCode,
			Link:         link,
		}
	} else {
		translationValues["signinsms"] = struct {
			Code string
			Link string
		}{
			Code: sessionInfo.SMSCode,
			Link: link,
		}
	}

	translations, err := tools.ParseTranslations(values.LangKey, translationValues)
	if err != nil {
		log.Error("Failed to parse translations: ", err)
		return
	}

	smsmessage := ""
	if authenticatingOrganization != "" {
		smsmessage = translations["authorizeorganizationsms"]
	} else {
		smsmessage = translations["signinsms"]
	}

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

// PhonenumberValidationAndLogin is the page that is linked to in the SMS for phonenumbervalidation
// and login. Therefore it is accessed on the mobile phone
func (service *Service) PhonenumberValidationAndLogin(w http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	values := request.Form
	key := values.Get("k")
	smscode := values.Get("c")
	langKey := values.Get("l")

	translationValues := tools.TranslationValues{
		"invalidlink":          nil,
		"error":                nil,
		"smsconfirmedandlogin": nil,
		"return_to_window":     nil,
	}

	translations, err := tools.ParseTranslations(langKey, translationValues)
	if err != nil {
		log.Error("Failed to parse translations: ", err)
		return
	}

	err = service.phonenumberValidationService.ConfirmValidation(request, key, smscode)
	if err == validation.ErrInvalidCode || err == validation.ErrInvalidOrExpiredKey {
		service.renderSMSConfirmationPage(w, request, translations["invalidlink"], "")
		return
	}
	if err != nil {
		log.Error(err)
		service.renderSMSConfirmationPage(w, request, translations["error"], "")
		return
	}

	sessionInfo, err := service.getLoginSessionInformation(request, key)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if sessionInfo == nil {
		service.renderSMSConfirmationPage(w, request, translations["invalidlink"], "")
		return
	}

	validsmscode := (smscode == sessionInfo.SMSCode)

	if !validsmscode { //TODO: limit to 3 failed attempts
		service.renderSMSConfirmationPage(w, request, translations["invalidlink"], "")
		return
	}
	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoLoginCollectionName)

	_, err = mgoCollection.UpdateAll(bson.M{"sessionkey": key}, bson.M{"$set": bson.M{"confirmed": true}})
	if err != nil {
		log.Error("Error while confirming sms 2fa - ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	service.renderSMSConfirmationPage(w, request, translations["smsconfirmedandlogin"], translations["return_to_window"])
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
	langKey := values.Get("l")

	translationValues := tools.TranslationValues{
		"invalidlink":      nil,
		"smsloggingin":     nil,
		"return_to_window": nil,
	}

	translations, err := tools.ParseTranslations(langKey, translationValues)
	if err != nil {
		log.Error("Failed to parse translations: ", err)
		return
	}

	var validsmscode bool
	sessionInfo, err := service.getLoginSessionInformation(request, sessionKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if sessionInfo == nil {
		service.renderSMSConfirmationPage(w, request, translations["invalidlink"], "")
		return
	}

	validsmscode = (smscode == sessionInfo.SMSCode)

	if !validsmscode { //TODO: limit to 3 failed attempts
		service.renderSMSConfirmationPage(w, request, translations["invalidlink"], "")
		return
	}
	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoLoginCollectionName)

	_, err = mgoCollection.UpdateAll(bson.M{"sessionkey": sessionKey}, bson.M{"$set": bson.M{"confirmed": true}})
	if err != nil {
		log.Error("Error while confirming sms 2fa - ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	service.renderSMSConfirmationPage(w, request, translations["smsloggingin"], translations["return_to_window"])
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

// LoginWithoutTwoFA gets user logged in without 2 FA
func (service *Service) LoginWithoutTwoFA(w http.ResponseWriter, request *http.Request) {
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

	_, err = service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error("Failed to get loginsession: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userMgr := user.NewManager(request)
	userMgr.RemoveExpireDate(username)

	//add last 2fa date if logging in with oauth2
	service.storeLast2FALogin(request, username)

	service.loginUser(w, request, username)
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
	if sessionInfo != nil && !sessionInfo.Confirmed {
		//Already confirmed on the phone
		validsmscode := (values.Smscode == sessionInfo.SMSCode)

		if !validsmscode {
			// TODO: limit to 3 failed attempts
			w.WriteHeader(422)
			log.Debugf("Expected code %s, got %s", sessionInfo.SMSCode, values.Smscode)
			return
		}
	}
	loginSession, err := service.GetSession(request, SessionLogin, "loginsession")
	if err != nil {
		log.Error("Failed to get loginsession: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	validationkey, _ := loginSession.Values["phonenumbervalidationkey"].(string)
	err = service.phonenumberValidationService.ConfirmValidation(request, validationkey, values.Smscode)
	if err == validation.ErrInvalidCode {
		log.Debug("Invalid code")
		// TODO: limit to 3 failed attempts
		w.WriteHeader(422)
		log.Debug("invalid code")
		return
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
	service.login(w, request, username)
}

func (service *Service) loginOauthUser(w http.ResponseWriter, request *http.Request, username string) {
	if err := service.SetLoggedInOauthUser(w, request, username); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sessions.Save(request, w)
	log.Debugf("Successfull oauth login without 2 factor authentication by '%s'", username)
	service.login(w, request, username)
}

func (service *Service) login(w http.ResponseWriter, request *http.Request, username string) {

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

	inviteCode := queryValues.Get("invitecode")
	if inviteCode != "" {
		err := verifyInfoAfterLogin(request, username, inviteCode)
		if err != nil {
			log.Error("Error while running verifyInfoAfterLogin: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}

	sessions.Save(request, w)
	response := struct {
		RedirectUrl string `json:"redirecturl"`
	}{}
	response.RedirectUrl = redirectURL
	log.Debug("Redirecting to:", redirectURL)
	json.NewEncoder(w).Encode(response)
}

func verifyInfoAfterLogin(request *http.Request, username string, inviteCode string) error {
	invitationMgr := invitations.NewInvitationManager(request)
	orgMgr := organizationdb.NewManager(request)
	valMgr := validationdb.NewManager(request)
	userMgr := user.NewManager(request)
	invite, err := invitationMgr.GetByCode(inviteCode)
	if err == mgo.ErrNotFound || invite.Status != invitations.RequestPending {
		// silently ignore
		return nil
	}
	if err != nil {
		return err
	}
	org, err := orgMgr.GetByName(invite.Organization)
	if org == nil {
		log.Warn("Cannot accept invitation of deleted organization: ", invite.Organization)
		return nil
	}
	if invite.Role == invitations.RoleMember {
		err = orgMgr.SaveMember(org, username)
		if err != nil {
			return err
		}
	} else if invite.Role == invitations.RoleOwner {
		err = orgMgr.SaveOwner(org, username)
		if err != nil {
			return err
		}
	}
	if invite.Method == invitations.MethodEmail {
		// Set this email address as verified and create a new one if necessary
		emailAddress := user.EmailAddress{Label: invite.EmailAddress, EmailAddress: invite.EmailAddress}
		err = userMgr.SaveEmail(username, emailAddress)
		if err != nil {
			return err
		}
		validatedEmailAddress := valMgr.NewValidatedEmailAddress(username, invite.EmailAddress)
		err = valMgr.SaveValidatedEmailAddress(validatedEmailAddress)
		if err != nil {
			return err
		}
	} else if invite.Method == invitations.MethodPhone {
		// Set this phone number as verified and create a new one if necessary
		phoneNumber := user.Phonenumber{Label: invite.PhoneNumber, Phonenumber: invite.PhoneNumber}
		err = userMgr.SavePhone(username, phoneNumber)
		if err != nil {
			return err
		}
		validatedPhoneNumber := valMgr.NewValidatedPhonenumber(username, invite.PhoneNumber)
		err = valMgr.SaveValidatedPhonenumber(validatedPhoneNumber)
		if err != nil {
			return err
		}
	}
	return invitationMgr.SetAcceptedByCode(inviteCode)
}

// ValidateEmail is the handler for POST /login/validateemail
func (service *Service) ValidateEmail(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Username string `json:"username"`
		LangKey  string `json:"langkey"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Debug("Error decoding the validte email request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// try to load the username if it is a validated phone number
	valMgr := validationdb.NewManager(r)
	userFromPhone, err := valMgr.GetByPhoneNumber(body.Username)
	if err != nil && !db.IsNotFound(err) {
		log.Error("Error while retrieving username from phone number: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !db.IsNotFound(err) {
		body.Username = userFromPhone.Username
	}

	userMgr := user.NewManager(r)
	user, err := userMgr.GetByName(body.Username)
	if err != nil {
		if db.IsNotFound(err) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		log.Error("Error while retrieving username: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(user.EmailAddresses) < 1 {
		log.Debug("User does not have any email addresses.")
		http.Error(w, http.StatusText(http.StatusPreconditionFailed), http.StatusPreconditionFailed)
		return
	}

	// Don't send verification if at least 1 email address is already verified
	ve, err := valMgr.GetByUsernameValidatedEmailAddress(body.Username)
	if err != nil && !db.IsNotFound(err) {
		log.Error("Error while retrieving verified email addresses: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//User has verified email addresses
	if err == nil && len(ve) > 0 {
		log.Debug("User has verified email addresses, rejecting validate request")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	// Don't send verification if one is already ongoing
	ov, err := valMgr.GetOngoingEmailAddressValidationByUser(body.Username)
	if err != nil && !db.IsNotFound(err) {
		log.Error("Error while checking ongoing email address verifications: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err == nil && len(ov) > 0 {
		log.Debug("User has ongoing email address verifications, rejecting validate request")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	for _, email := range user.EmailAddresses {
		_, err = service.emailaddressValidationService.RequestValidation(r, body.Username, email.EmailAddress, fmt.Sprintf("https://%s/emailvalidation", r.Host), body.LangKey)
		if err != nil {
			log.Error("Failed to validate email address: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

//ForgotPassword handler for POST /login/forgotpassword
func (service *Service) ForgotPassword(w http.ResponseWriter, request *http.Request) {
	// login can be username or email
	values := struct {
		Login   string `json:"login"`
		LangKey string `json:"langkey"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&values); err != nil {
		log.Debug("Error decoding the ForgotPassword request:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Make field case insensitive
	values.Login = strings.ToLower(values.Login)
	userMgr := user.NewManager(request)
	valMgr := validationdb.NewManager(request)
	validatedemail, err := valMgr.GetByEmailAddressValidatedEmailAddress(values.Login)
	if err != nil && !db.IsNotFound(err) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	var username string
	var emails []string
	// If the login value is a validated email address, use that to get the username
	// and only send the reset email to this email address
	if !db.IsNotFound(err) {
		username = validatedemail.Username
		emails = []string{validatedemail.EmailAddress}
	} else {
		// First check if the login value is a phone number, and update the login
		// value with the username coupled to the phone number if it is. Then just get
		// all validated email addresses for the user and send emails to them
		validatedPhone, err := valMgr.GetByPhoneNumber(values.Login)
		if err != nil && !db.IsNotFound(err) {
			log.Error("Failed to get validated phone numbers: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		if !db.IsNotFound(err) {
			values.Login = validatedPhone.Username
		}
		usr, err := userMgr.GetByName(values.Login)
		if err != nil && err != mgo.ErrNotFound || usr.Username == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		username = usr.Username
		validatedemails, err := valMgr.GetByUsernameValidatedEmailAddress(username)
		if validatedemails == nil || len(validatedemails) == 0 {
			http.Error(w, http.StatusText(http.StatusPreconditionFailed), http.StatusPreconditionFailed)
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
	_, err = service.emailaddressValidationService.RequestPasswordReset(request, username, emails, values.LangKey)
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
		LangKey     string `json:"langkey"`
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
	if !phonenumber.Validate() {
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

	// save the phone number validation request
	valMngr := validationdb.NewManager(request)
	info, err := valMngr.NewPhonenumberValidationInformation(username, phonenumber)
	if err != nil {
		log.Error("ResendPhonenumberConfirmation: Could not create phonenumber validation information: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = valMngr.SavePhonenumberValidationInformation(info)
	if err != nil {
		log.Error("ResendPhonenumberConfirmation: Could not save phonenumber validation information: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	loginSession.Values["phonenumbervalidationkey"] = info.Key

	sessionInfo := &loginSessionInformation{
		Confirmed:  false,
		CreatedAt:  time.Now(),
		SessionKey: info.Key,
		SMSCode:    info.SMSCode,
	}

	loginSession.Values["sessionkey"] = sessionInfo.SessionKey

	mgoCollection := db.GetCollection(db.GetDBSession(request), mongoLoginCollectionName)
	mgoCollection.Insert(sessionInfo)

	TranslationValues := tools.TranslationValues{
		"smsconfirmationandlogin": struct {
			Code string
			Link string
		}{
			Code: info.SMSCode,
			Link: fmt.Sprintf("https://%s/pvl?c=%s&k=%s&l=%s", request.Host, info.SMSCode, url.QueryEscape(info.Key), values.LangKey),
		},
	}

	translations, err := tools.ParseTranslations(values.LangKey, TranslationValues)
	if err != nil {
		log.Error("Failed to parse translations: ", err)
		return
	}

	go service.phonenumberValidationService.SMSService.Send(values.PhoneNumber,
		translations["smsconfirmationandlogin"])

	sessions.Save(request, w)
	w.WriteHeader(http.StatusNoContent)
}

func (service *Service) GetOrganizationInvitation(w http.ResponseWriter, request *http.Request) {
	code := mux.Vars(request)["code"]
	if code == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	invitationMgr := invitations.NewInvitationManager(request)
	invite, err := invitationMgr.GetByCode(code)
	if err == mgo.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(invite)
	}
}
