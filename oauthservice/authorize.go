package oauthservice

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/credentials/oauth2"
	organizationdb "github.com/itsyouonline/identityserver/db/organization"
)

type authorizationRequest struct {
	AuthorizationCode string
	Username          string
	RedirectURL       string
	ClientID          string
	State             string
	Scope             string
	CreatedAt         time.Time
}

func (ar *authorizationRequest) IsExpiredAt(testtime time.Time) bool {
	return testtime.After(ar.CreatedAt.Add(time.Second * 10))
}

func newAuthorizationRequest(username, clientID, state, scope, redirectURI string) *authorizationRequest {
	var ar authorizationRequest
	randombytes := make([]byte, 21) //Multiple of 3 to make sure no padding is added
	rand.Read(randombytes)
	ar.AuthorizationCode = base64.URLEncoding.EncodeToString(randombytes)
	ar.CreatedAt = time.Now()
	ar.Username = username
	ar.ClientID = clientID
	ar.State = state
	ar.Scope = scope
	ar.RedirectURL = redirectURI

	return &ar
}

func validateRedirectURI(mgr ClientManager, redirectURI string, clientID string) (valid bool, err error) {
	log.Debug("Validating redirect URI for ", clientID)
	u, err := url.Parse(redirectURI)
	if err != nil {
		err = nil
		return
	}

	valid = true
	//A redirect to itsyou.online can not do harm but it is not normal either
	valid = valid && (u.Scheme != "")
	lowercaseHost := strings.ToLower(u.Host)
	valid = valid && (lowercaseHost != "")
	valid = valid && (!strings.HasSuffix(lowercaseHost, "itsyou.online"))
	valid = valid && (!strings.Contains(lowercaseHost, "itsyou.online:"))

	if !valid {
		return
	}

	//For now, just check if the redirectURI is registered in 'a' apikey
	//The redirect_uri is saved in the authorization request and during
	// the access_token request when the secret is available, check again against the known value
	clients, err := mgr.AllByClientID(clientID)
	if err != nil {
		valid = false
		return
	}

	match := false
	for _, client := range clients {
		log.Debug("Possible redirect_uri: ", client.Label, "\n ", client.CallbackURL)
		match = match || strings.HasPrefix(redirectURI, client.CallbackURL)
	}
	valid = valid && match

	log.Debug("Redirect URI is valid: ", valid)
	return
}

func redirectToNextPage(w http.ResponseWriter, r *http.Request) {
	queryvalues := r.URL.Query()
	queryvalues.Add("endpoint", r.URL.EscapedPath())
	redirectToRegistrationPage := r.Form.Get("register") != ""
	//TODO: redirect according the the received http method
	if redirectToRegistrationPage {
		http.Redirect(w, r, "/register?"+queryvalues.Encode(), http.StatusFound)
	} else {
		http.Redirect(w, r, "/login?"+queryvalues.Encode(), http.StatusFound)
	}
}

func redirectToScopeRequestPage(w http.ResponseWriter, r *http.Request, possibleScopes []string) {
	var possibleScopesString string
	if possibleScopes != nil {
		possibleScopesString = strings.Join(possibleScopes, ",")
	}
	queryvalues := r.URL.Query()
	queryvalues.Set("scope", possibleScopesString)
	queryvalues.Add("endpoint", r.URL.EscapedPath())
	//TODO: redirect according the the received http method
	http.Redirect(w, r, "/authorize?"+queryvalues.Encode(), http.StatusFound)
}

func (service *Service) filterAuthorizedScopes(r *http.Request, username string, clientID string, requestedScopes []string) (authorizedScopes []string, err error) {
	log.Debug("Validating authorizations for requested scopes: ", requestedScopes)
	authorizedScopes, err = service.identityService.FilterAuthorizedScopes(r, username, clientID, requestedScopes)
	log.Debug("Authorized scopes: ", authorizedScopes)
	//TODO: how to request explicit confirmation?

	return
}

//AuthorizeHandler is the handler of the /v1/oauth/authorize endpoint
func (service *Service) AuthorizeHandler(w http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		log.Debug("ERROR parsing form", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Check if the requested authorization grant type is supported
	requestedResponseType := request.Form.Get("response_type")
	if requestedResponseType != AuthorizationGrantCodeType {
		log.Debug("Invalid authorization grant type requested")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Check if the user is already authenticated, if not, redirect to the login page before returning here
	var protectedSession bool
	username, err := service.GetWebuser(request, w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if username == "" {
		username, err = service.GetOauthUser(request, w)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if username != "" {
			log.Debug("protected session")
			protectedSession = true
		} else {
			redirectToNextPage(w, request)
			return
		}
	}

	//Validate client and redirect_uri
	redirectURI, err := url.QueryUnescape(request.Form.Get("redirect_uri"))
	if err != nil {
		log.Debug("Unparsable redirect_uri")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	clientID := request.Form.Get("client_id")
	mgr := NewManager(request)
	valid, err := validateRedirectURI(mgr, redirectURI, clientID)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	requestedScopes := oauth2.SplitScopeString(request.Form.Get("scope"))
	possibleScopes, err := service.filterPossibleScopes(request, username, requestedScopes, true)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	authorizedScopes, err := service.filterAuthorizedScopes(request, username, clientID, possibleScopes)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var authorizedScopeString string
	var validAuthorization bool

	if authorizedScopes != nil {
		authorizedScopeString = strings.Join(authorizedScopes, ",")
		validAuthorization = len(possibleScopes) == len(authorizedScopes)
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

	//If no valid authorization, ask the user for authorizations
	if !validAuthorization {
		if protectedSession {
			log.Debug("protected session active, but need to give authorizations")
			// We need a full session to give authorizations, so remove the l2fa entry
			// This way the login function will require 2fa and give a full session with admin scopes
			l2faMgr := organizationdb.NewLast2FAManager(request)
			if l2faMgr.Exists(clientID, username) {
				err = l2faMgr.RemoveLast2FA(clientID, username)
				if err != nil {
					log.Error(err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}
			redirectToNextPage(w, request)
			return
		}
		token, e := service.createItsYouOnlineAdminToken(username, request)
		if e != nil {
			log.Error(e)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		service.sessionService.SetAPIAccessToken(w, token)
		redirectToScopeRequestPage(w, request, possibleScopes)
		return
	}

	if clientID == "itsyouonline" {
		log.Warn("HACK attempt, someone tried to get a token as the 'itsyouonline' client")
		//TODO: log the entire request and everything we know
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	redirectURI, err = handleAuthorizationGrantCodeType(request, username, clientID, redirectURI, authorizedScopeString)

	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, request, redirectURI, http.StatusFound)

}

func handleAuthorizationGrantCodeType(r *http.Request, username, clientID, redirectURI, scopes string) (correctedRedirectURI string, err error) {
	correctedRedirectURI = redirectURI
	log.Debug("Handling authorization grant code type for user ", username, ", ", clientID, " is asking for ", scopes)
	clientState := r.Form.Get("state")
	//TODO: validate state (length and stuff)

	ar := newAuthorizationRequest(username, clientID, clientState, scopes, redirectURI)
	mgr := NewManager(r)
	err = mgr.saveAuthorizationRequest(ar)
	if err != nil {
		return
	}

	parameters := make(url.Values)
	parameters.Add("code", ar.AuthorizationCode)
	parameters.Add("state", clientState)

	//Don't parse the redirect url, can only give errors while we don't gain much
	if !strings.Contains(correctedRedirectURI, "?") {
		correctedRedirectURI += "?"
	} else {
		if !strings.HasSuffix(correctedRedirectURI, "&") {
			correctedRedirectURI += "&"
		}
	}
	correctedRedirectURI += parameters.Encode()
	return
}
