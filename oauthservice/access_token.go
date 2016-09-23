package oauthservice

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/db/user/apikey"
)

//AccessTokenExpiration is the time in seconds an access token expires
var AccessTokenExpiration = time.Second * 3600 * 24 //Tokens expire after 1 day

//AccessToken is an oauth2 accesstoken together with the access information it stands for
type AccessToken struct {
	AccessToken string
	Type        string
	Username    string
	GlobalID    string //The organization that granted the token (in case of a client credentials flow)
	Scope       string
	ClientID    string //The client_id of the organization that was granted the token
	CreatedAt   time.Time
}

//IsExpiredAt checks if the token is expired at a specific time
func (at *AccessToken) IsExpiredAt(testtime time.Time) bool {
	return testtime.After(at.ExpirationTime())
}

//IsExpired is a convenience method for IsExpired(time.Now())
func (at *AccessToken) IsExpired() bool {
	return at.IsExpiredAt(time.Now())
}

//ExpirationTime return the time at which this token expires
func (at *AccessToken) ExpirationTime() time.Time {
	return at.CreatedAt.Add(AccessTokenExpiration)
}

func newAccessToken(username, globalID, clientID, scope string) *AccessToken {
	var at AccessToken

	randombytes := make([]byte, 21) //Multiple of 3 to make sure no padding is added
	rand.Read(randombytes)
	at.AccessToken = base64.URLEncoding.EncodeToString(randombytes)
	at.CreatedAt = time.Now()
	at.Username = username
	at.GlobalID = globalID
	at.ClientID = clientID
	at.Scope = scope
	at.Type = "bearer"

	return &at
}

//AccessTokenHandler is the handler of the /v1/oauth/access_token endpoint
func (service *Service) AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing form: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	grantType := r.FormValue("grant_type")
	clientSecret := r.FormValue("client_secret")
	clientID := r.FormValue("client_id")

	//Also accept some alternatives
	if grantType == "authorization_code" {
		grantType = ""
	}

	if clientSecret == "" || clientID == "" || (grantType == "" && code == "") {
		log.Debug("Required parameter missing in the request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var at *AccessToken
	httpStatusCode := http.StatusOK

	mgr := NewManager(r)
	if grantType != "" {
		if grantType == ClientCredentialsGrantCodeType {
			at, httpStatusCode = clientCredentialsTokenHandler(clientID, clientSecret, mgr, r)
		} else {
			log.Debug("Invalid grant_type")
			httpStatusCode = http.StatusBadRequest
		}
	} else {
		redirectURI := r.FormValue("redirect_uri")
		at, httpStatusCode = convertCodeToAccessTokenHandler(code, clientID, clientSecret, redirectURI, mgr, r)
	}

	if httpStatusCode != http.StatusOK {
		http.Error(w, http.StatusText(httpStatusCode), httpStatusCode)
		return
	}

	//It is also possible to immediately get a JWT by specifying 'id_token' as the response type
	// In this case, the scope parameter needs to be given to prevent consumers to accidentially handing out too powerful tokens to third party services
	// It is also possible to specify additional audiences
	responseType := r.FormValue("response_type")

	if responseType == "id_token" {
		requestedScopeParameter := r.FormValue("scope")
		extraAudiences := r.FormValue("aud")
		tokenString, err := service.convertAccessTokenToJWT(r, at, requestedScopeParameter, extraAudiences)
		if err == errUnauthorized {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/jwt")
		w.Write([]byte(tokenString))
		return
	}
	mgr.saveAccessToken(at)

	response := struct {
		AccessToken string      `json:"access_token"`
		TokenType   string      `json:"token_type"`
		Scope       string      `json:"scope"`
		ExpiresIn   int64       `json:"expires_in"`
		Info        interface{} `json:"info"`
	}{
		AccessToken: at.AccessToken,
		TokenType:   at.Type,
		Scope:       at.Scope,
		ExpiresIn:   int64(at.ExpirationTime().Sub(time.Now()).Seconds() - 600),

		Info: struct {
			Username string `json:"username"`
		}{
			Username: at.Username,
		},
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

func clientCredentialsTokenHandler(clientID string, secret string, mgr *Manager, r *http.Request) (at *AccessToken, httpStatusCode int) {
	httpStatusCode = http.StatusOK
	var scopes string
	username := ""

	client, err := mgr.getClientByCredentials(clientID, secret)
	if err != nil {
		log.Error("Error getting the oauth client: ", err)
		httpStatusCode = http.StatusInternalServerError
		return
	}
	if client == nil || !client.ClientCredentialsGrantType {
		log.Info("Checking user api")
		apikeyMgr := apikey.NewManager(r)
		apikey, err := apikeyMgr.GetByApplicationAndSecret(clientID, secret)
		if err != nil || apikey.ApiKey != secret {
			log.Error("Error getting the user api key: ", err)
			httpStatusCode = http.StatusBadRequest
			return
		}
		log.Info("apikey", apikey)
		scopes = strings.Join(apikey.Scopes, " ")
		log.Info("scopes ", scopes)
		username = apikey.Username
	} else {
		scopes = "organization:owner"
	}

	at = newAccessToken(username, clientID, clientID, scopes)
	return
}

func convertCodeToAccessTokenHandler(code string, clientID string, secret string, redirectURI string, mgr *Manager, r *http.Request) (at *AccessToken, httpStatusCode int) {
	httpStatusCode = http.StatusOK

	ar, err := mgr.Get(code)
	if err != nil {
		log.Error("ERROR getting the original authorization request:", err)
		httpStatusCode = http.StatusInternalServerError
		return
	}
	if ar == nil {
		log.Debug("No original authorization request found with this authorization code")
		httpStatusCode = http.StatusBadRequest
		return
	}

	state := r.FormValue("state")

	if ar.ClientID != clientID || ar.State != state || ar.RedirectURL != redirectURI {
		log.Info("Bad client or hacking attempt, state, client_id or redirect_uri is different from the original authorization request")
		httpStatusCode = http.StatusBadRequest
		return
	}

	if ar.IsExpiredAt(time.Now()) {
		log.Info("Token request for an expired authorizationrequest")
		httpStatusCode = http.StatusBadRequest
		return
	}

	client, err := mgr.getClientByCredentials(clientID, secret)
	if err != nil {
		log.Error("Error getting the oauth client: ", err)
		httpStatusCode = http.StatusInternalServerError
		return
	}
	if client == nil {
		log.Info("(client_id - secret) combination not found")
		httpStatusCode = http.StatusBadRequest
		return
	}

	if !strings.HasPrefix(redirectURI, client.CallbackURL) {
		log.Debug("return_uri does not match the callback uri")
		httpStatusCode = http.StatusBadRequest
		return
	}

	at = newAccessToken(ar.Username, "", ar.ClientID, ar.Scope)
	return
}

func (service *Service) createItsYouOnlineAdminToken(username string, r *http.Request) (token string, err error) {
	at := newAccessToken(username, "", "itsyouonline", "admin")

	mgr := NewManager(r)
	err = mgr.saveAccessToken(at)
	if err == nil {
		token = at.AccessToken
	}
	return
}
