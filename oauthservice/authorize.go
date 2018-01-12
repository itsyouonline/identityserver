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
	"github.com/itsyouonline/identityserver/db"
	organizationdb "github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/db/validation"
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
		validAuthorization = IsAuthorizationValid(possibleScopes, authorizedScopes)

		// Check if the user still has the given authorizations
		authorization, err := user.NewManager(request).GetAuthorization(username, clientID)
		if err != nil {
			log.Error("Failed to load authorization: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if validAuthorization {
			validAuthorization, err = UserHasAuthorizedScopes(request, authorization)
			if err != nil {
				log.Error("Failed to check if authorizated labels are still present: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
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
	log.Debug("Redirecting from authorize handler to: ", redirectURI)
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

// IsAuthorizationValid checks if the possible scopes that are being requested are already authorized
func IsAuthorizationValid(possibleScopes []string, authorizedScopes []string) bool {
	if len(possibleScopes) > len(authorizedScopes) {
		return false
	}
POSSIBLESCOPES:
	for _, possibleScope := range possibleScopes {
		for _, authorizedScope := range authorizedScopes {
			if possibleScope == authorizedScope {
				// Scope is already authorized, move on to the next one
				continue POSSIBLESCOPES
			}
		}
		// If we get here, it means the possibleScope is not in the authorizedScopes list
		// So the standing authorization is not valid
		return false
	}
	// Likewise, if we get here we did not yet return due to a missing scope authorization
	// and all possible scopes are exhausted. Therefore the standing authorization
	// is valid
	return true
}

// UserHasAuthorizedScopes checks if all labels from an authorization scope mapping are still present on the user
func UserHasAuthorizedScopes(r *http.Request, authorization *user.Authorization) (bool, error) {
	log.Debug("Checking if user still has all mapped scopes")
	user, err := user.NewManager(r).GetByName(authorization.Username)
	if err != nil {
		return false, err
	}

	found := false
	if authorization.Addresses != nil {
		for _, am := range authorization.Addresses {
			for _, a := range user.Addresses {
				if am.RealLabel == a.Label {
					found = true
					break
				}
			}
			if !found {
				log.Debug("Authorized real label not found for address")
				return false, nil
			}
		}
	}

	found = false
	if authorization.BankAccounts != nil {
		for _, ba := range authorization.BankAccounts {
			for _, b := range user.BankAccounts {
				if ba.RealLabel == b.Label {
					found = true
					break
				}
			}
			if !found {
				log.Debug("Authorized real label not found for bank account")
				return false, nil
			}
		}
	}

	found = false
	if authorization.DigitalWallet != nil {
		for _, dw := range authorization.DigitalWallet {
			for _, d := range user.DigitalWallet {
				if dw.RealLabel == d.Label {
					found = true
					break
				}
			}
			if !found {
				log.Debug("Authorized real label not found for digital wallet")
				return false, nil
			}
		}
	}

	found = false
	if authorization.EmailAddresses != nil {
		for _, ea := range authorization.EmailAddresses {
			for _, e := range user.EmailAddresses {
				if ea.RealLabel == e.Label {
					found = true
					break
				}
			}
			if !found {
				log.Debug("Authorized real label not found for email address")
				return false, nil
			}
		}
	}

	found = false
	if authorization.Phonenumbers != nil {
		for _, ep := range authorization.Phonenumbers {
			for _, p := range user.Phonenumbers {
				if ep.RealLabel == p.Label {
					found = true
					break
				}
			}
			if !found {
				log.Debug("Authorized real label not found for phone number")
				return false, nil
			}
		}
	}

	found = false
	if authorization.PublicKeys != nil {
		for _, pk := range authorization.PublicKeys {
			for _, p := range user.PublicKeys {
				if pk.RealLabel == p.Label {
					found = true
					break
				}
			}
			if !found {
				log.Debug("Authorized real label not found for public key")
				return false, nil
			}
		}
	}

	found = false
	valMgr := validation.NewManager(r)
	if authorization.ValidatedEmailAddresses != nil {
		for _, vea := range authorization.ValidatedEmailAddresses {
			for _, e := range user.EmailAddresses {
				if vea.RealLabel == e.Label {
					// Check that the email address is actually validated by the user in the authorization
					validatedEmail, err := valMgr.GetByEmailAddressValidatedEmailAddress(e.EmailAddress)
					if err != nil && !db.IsNotFound(err) {
						return false, err
					}
					if db.IsNotFound(err) {
						return false, nil
					}
					if validatedEmail.Username == authorization.Username {
						found = true
						break
					}
				}
			}
			if !found {
				log.Debug("Validated email address in authorization does not belong to the user")
				return false, nil
			}
		}
	}

	found = false
	if authorization.ValidatedPhonenumbers != nil {
		for _, vpn := range authorization.ValidatedPhonenumbers {
			for _, p := range user.Phonenumbers {
				if vpn.RealLabel == p.Label {
					// Check that the phone is actually validated by the user in the authorization
					validatedPhone, err := valMgr.GetByPhoneNumber(p.Phonenumber)
					if err != nil && !db.IsNotFound(err) {
						return false, err
					}
					if db.IsNotFound(err) {
						return false, nil
					}
					if validatedPhone.Username == authorization.Username {
						found = true
						break
					}
				}
			}
			if !found {
				log.Debug("Validated phone number in authorization does not belong to the user")
				return false, nil
			}
		}
	}
	return true, nil
}
