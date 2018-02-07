package oauthservice

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/itsyouonline/identityserver/credentials/oauth2"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/grants"
	"github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/db/validation"
)

var errUnauthorized = errors.New("Unauthorized")

const issuer = "itsyouonline"

//JWTHandler returns a JWT with claims that are a subset of the scopes available to the authorizing token
func (service *Service) JWTHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Debug("Error parsing form: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	requestedScopeParameter := r.FormValue("scope")
	audiences := strings.TrimSpace(r.FormValue("aud"))

	//First check if the user uses an existing jwt to authenticate and authorize itself
	idToken, err := oauth2.GetValidJWT(r, &service.jwtSigningKey.PublicKey)
	if err != nil {
		log.Warning(err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var tokenString string
	if idToken != nil {
		tokenString, err = service.createNewJWTFromParent(r, idToken, requestedScopeParameter, audiences)
	} else {
		//If no jwt was supplied, check if an old school access_token was used
		accessToken := r.Header.Get("Authorization")

		//Get the actual token out of the header (accept 'token ABCD' as well as just 'ABCD' and ignore some possible whitespace)
		accessToken = strings.TrimSpace(strings.TrimPrefix(accessToken, "token"))
		if accessToken == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		oauthMgr := NewManager(r)
		var at *AccessToken
		at, err = oauthMgr.GetAccessToken(accessToken)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if at == nil || at.IsExpired() {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		validity := parseValidity(r)

		tokenString, err = service.convertAccessTokenToJWT(r, at, requestedScopeParameter, audiences, validity)
	}
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
}

//RefreshJWTHandler returns a new refreshed JWT with the same scopes as the original JWT
// The original JWT needs to be passed in the authorization header as a bearer token
// If the stored allowed scopes no longer contains a specific scope present in the jwt, this scope is also dropped in the newly created JWT.
func (service *Service) RefreshJWTHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Debug("Error parsing form: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	originalToken, err := oauth2.GetValidJWT(r, &service.jwtSigningKey.PublicKey)
	err = oauth2.IgnoreExpired(err)
	if err != nil {
		log.Warning(err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if originalToken == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	mgr := NewManager(r)
	rawRefreshToken, refreshtokenPresent := originalToken.Claims["refresh_token"]
	if !refreshtokenPresent {
		log.Debug("No refresh_token in the jwt supplied:", originalToken)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	refreshTokenString, ok := rawRefreshToken.(string)
	if !ok {
		log.Error("ERROR while reading the refresh token from the jwt")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	rt, err := mgr.getRefreshToken(refreshTokenString)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rt == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// Take the scope from the stored refreshtoken, it might be that certain authorizations are revoked
	// Also validate a possible memberof:clientId scope
	orgMgr := organization.NewManager(r)
	clientID := originalToken.Claims["azp"].(string)
	username, isUser := originalToken.Claims["username"].(string)
	// if a username is set verify the possible membership scopes.
	scope := strings.Join(rt.Scopes, ",")
	if isUser {
		scope, err = verifyScopes(scope, username, clientID, orgMgr)
		if err != nil {
			log.Error("Error while verifying scopes for user jwt: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	originalToken.Claims["scope"] = strings.Split(scope, ",")

	// Add grants
	if r.FormValue("add_grants") == "true" && isUser {
		grantList, err := getGrants(username, clientID, r)
		if err != nil {
			log.Error("Failed to add grants in jwt refresh: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		originalToken.Claims["scope"] = append(strings.Split(scope, ","), grantList...)
	}

	// Set a new expiration time
	validity := parseValidity(r)

	expiration := time.Now().Add(AccessTokenExpiration).Unix()

	requestedExpiration := expiration
	if validity > 0 {
		requestedExpiration = time.Now().Unix() + validity
		if requestedExpiration < expiration {
			expiration = requestedExpiration
		}
	}
	originalToken.Claims["exp"] = expiration
	// Sign it and return
	tokenString, err := originalToken.SignedString(service.jwtSigningKey)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rt.LastUsed = db.DateTime(time.Now())
	err = mgr.saveRefreshToken(rt)
	if err != nil {
		log.Error("Error while saving refresh token:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/jwt")
	w.Write([]byte(tokenString))
}

func stripOfflineAccess(scopes []string) (result []string, offlineAccessRequested bool) {
	result = make([]string, 0, len(scopes))
	for _, scope := range scopes {
		if scope == "offline_access" {
			offlineAccessRequested = true
		} else {
			result = append(result, scope)
		}
	}
	return
}

func (service *Service) convertAccessTokenToJWT(r *http.Request, at *AccessToken, requestedScopeString, audiences string, maxValid int64) (tokenString string, err error) {
	requestedScopes := oauth2.SplitScopeString(requestedScopeString)
	requestedScopes, offlineAccessRequested := stripOfflineAccess(requestedScopes)
	acquiredScopes := oauth2.SplitScopeString(at.Scope)

	if len(requestedScopes) == 0 {
		// if the scope parameter is ommited, give all the authorized scopes
		// offline_access is already removed here, so just requesting that scope will
		// also give all requested scopes
		requestedScopes = acquiredScopes
	}

	//Basic validation to check if the requested scopes are possible within the acquiredScopes
	if !jwtScopesAreAllowed(acquiredScopes, requestedScopes) {
		err = errUnauthorized
		return
	}

	token := jwt.New(jwt.SigningMethodES384)

	//More extensive validation
	var grantedScopes []string
	if at.Username != "" {
		token.Claims["username"] = at.Username
		grantedScopes, err = service.filterPossibleScopes(r, at.Username, requestedScopes, false)
		if err != nil {
			return
		}
	}
	if at.GlobalID != "" {
		token.Claims["globalid"] = at.GlobalID
		grantedScopes = requestedScopes
	}
	// now that we have the granted scopes, check if the actual data is requested
	if r.FormValue("store_info") == "true" && at.Username != "" {
		grantedScopes = storeActualValue(r, grantedScopes, at.Username, at.ClientID)
	}
	token.Claims["scope"] = grantedScopes

	// process the audience string and make sure we don't set an empty slice if no
	// audience is set explicitly
	var audiencesArr []string
	for _, aud := range strings.Split(audiences, ",") {
		trimmedAud := strings.TrimSpace(aud)
		if trimmedAud != "" {
			audiencesArr = append(audiencesArr, trimmedAud)
		}
	}
	if len(audiencesArr) > 0 {
		token.Claims["aud"] = audiencesArr
	}

	// It does not hurt to always set the azp claim while it is only needed when the ID Token has a single
	// audience value and that audience is different than the authorized party
	token.Claims["azp"] = at.ClientID

	expiration := at.ExpirationTime().Unix()
	// If a custom validity period for the jwt is set, verify that it expires sooner
	// then the access_token would, and set that timestamp. If not, just keep the old expiration
	// timestamp
	requestedExpiration := expiration
	if maxValid > 0 {
		requestedExpiration = time.Now().Unix() + maxValid
		if requestedExpiration < expiration {
			expiration = requestedExpiration
		}
	}
	token.Claims["exp"] = expiration
	token.Claims["iss"] = issuer

	if offlineAccessRequested {
		rt := newRefreshToken()
		rt.AuthorizedParty = at.ClientID
		rt.Scopes = grantedScopes
		rt.LastUsed = db.DateTime(time.Now())
		token.Claims["refresh_token"] = rt.RefreshToken
		mgr := NewManager(r)
		if err = mgr.saveRefreshToken(&rt); err != nil {
			return
		}
	}
	orgMgr := organization.NewManager(r)
	scope, err := verifyScopes(strings.Join(grantedScopes, ","), at.Username, at.ClientID, orgMgr)
	if err != nil {
		return
	}
	token.Claims["scope"] = strings.Split(scope, ",")

	// Add grants
	if r.FormValue("add_grants") == "true" && at.Username != "" {
		grantList, err := getGrants(at.Username, at.ClientID, r)
		if err != nil {
			return "", err
		}
		token.Claims["scope"] = append(strings.Split(scope, ","), grantList...)
	}

	tokenString, err = token.SignedString(service.jwtSigningKey)
	return
}

func (service *Service) createNewJWTFromParent(r *http.Request, parentToken *jwt.Token, requestedScopeString, audiences string) (tokenString string, err error) {

	requestedScopes := oauth2.SplitScopeString(requestedScopeString)
	requestedScopes, offlineAccessRequested := stripOfflineAccess(requestedScopes)

	acquiredScopes := oauth2.GetScopesFromJWT(parentToken)
	var parentRefreshToken *refreshToken
	mgr := NewManager(r)
	if rawParentRefreshToken, parentRefreshTokenSupplied := parentToken.Claims["refresh_token"]; parentRefreshTokenSupplied {
		parentRefreshTokenString := rawParentRefreshToken.(string)
		parentRefreshToken, err = mgr.getRefreshToken(parentRefreshTokenString)
		if err != nil {
			return
		}
		if parentRefreshToken == nil {
			err = errUnauthorized
			return
		}
		acquiredScopes = parentRefreshToken.Scopes
	} else {
		// Do not allow a refreshtoken using a parent that does not have one
		if offlineAccessRequested {
			err = errUnauthorized
			return
		}
		//TODO: check if the parent token expired
	}

	if !jwtScopesAreAllowed(acquiredScopes, requestedScopes) {
		err = errUnauthorized
		return
	}

	token := jwt.New(jwt.SigningMethodES384)
	var grantedScopes []string
	username := parentToken.Claims["username"].(string)
	if username != "" {
		token.Claims["username"] = username
		grantedScopes, err = service.filterPossibleScopes(r, username, requestedScopes, false)
		if err != nil {
			return
		}
	}
	globalID := parentToken.Claims["globalid"]
	if globalID != "" {
		token.Claims["globalid"] = globalID
		grantedScopes = requestedScopes
	}
	if r.FormValue("store_info") == "true" && username != "" {
		grantedScopes = storeActualValue(r, grantedScopes, username, parentToken.Claims["azp"].(string))
	}
	token.Claims["scope"] = grantedScopes

	clientID, _ := parentToken.Claims["azp"].(string)
	if r.FormValue("add_grants") == "true" && username != "" {
		grantList, err := getGrants(username, clientID, r)
		if err != nil {
			return "", err
		}
		token.Claims["scope"] = append(grantedScopes, grantList...)
	}

	// process the audience string and make sure we don't set an empty slice if no
	// audience is set explicitly
	var audiencesArr []string
	for _, aud := range strings.Split(audiences, ",") {
		trimmedAud := strings.TrimSpace(aud)
		if trimmedAud != "" {
			audiencesArr = append(audiencesArr, trimmedAud)
		}
	}
	if len(audiencesArr) > 0 {
		token.Claims["aud"] = audiencesArr
	}
	token.Claims["azp"] = parentToken.Claims["azp"]
	lastUsed := db.DateTime(time.Now())
	if parentRefreshToken != nil {
		token.Claims["exp"] = time.Now().Add(AccessTokenExpiration).Unix()
		parentRefreshToken.LastUsed = lastUsed
		if err = mgr.saveRefreshToken(parentRefreshToken); err != nil {
			return
		}
	} else {
		token.Claims["exp"] = parentToken.Claims["exp"]
	}
	validity := parseValidity(r)
	expiration := token.Claims["exp"].(int64)
	requestedExpiration := expiration
	if validity > 0 {
		requestedExpiration = time.Now().Unix() + validity
		if requestedExpiration < expiration {
			expiration = requestedExpiration
		}
	}
	token.Claims["exp"] = expiration

	token.Claims["iss"] = issuer

	if offlineAccessRequested {
		rt := newRefreshToken()
		rt.Parent = parentRefreshToken.RefreshToken
		rt.AuthorizedParty = token.Claims["azp"].(string)
		rt.Scopes = grantedScopes
		rt.LastUsed = lastUsed
		token.Claims["refresh_token"] = rt.RefreshToken
		if err = mgr.saveRefreshToken(&rt); err != nil {
			return
		}
	}
	tokenString, err = token.SignedString(service.jwtSigningKey)
	return
}

func jwtScopesAreAllowed(grantedScopes []string, requestedScopes []string) (valid bool) {
	valid = true
	for _, rs := range requestedScopes {
		log.Debug(fmt.Sprintf("Checking if '%s' is allowed", rs))
		valid = valid && checkIfScopeInList(grantedScopes, rs)
	}

	return
}

func checkIfScopeInList(grantedScopes []string, scope string) (valid bool) {
	for _, as := range grantedScopes {
		//Allow all user scopes if the 'user:admin' scope is part of the autorized scopes
		if as == "user:admin" {
			if strings.HasPrefix(scope, "user:") {
				valid = true
				return
			}
		}
		if strings.HasPrefix(scope, as) {
			valid = true
			return
		}
	}
	return
}

// parseValidity parses the validity parameter from the request
func parseValidity(r *http.Request) int64 {
	validityString := r.FormValue("validity")
	var validity int64
	if validityString == "" {
		validity = -1
	} else {
		var err error
		validity, err = strconv.ParseInt(validityString, 10, 64)
		if err != nil {
			log.Debugf("Failed to parse validity argument (%v) as int64", validityString)
			validity = -1
		}
	}
	return validity
}

func storeActualValue(r *http.Request, grantedScopes []string, username string, clientID string) []string {
	if username == "" {
		return grantedScopes
	}

	fullScopes := []string{}
	for _, scope := range grantedScopes {
		fullScopes = append(fullScopes, scope)
		if strings.HasPrefix(scope, "user:phone") {
		}
	}
	clientCredentialsJWT := false
	orgMgr := organization.NewManager(r)
	if _, err := orgMgr.GetByName(clientID); db.IsNotFound(err) {
		clientCredentialsJWT = true
	}
	var authorization *user.Authorization
	// jwt aquired through code grant flow - there must be an authorization to represent
	// the scope mappings
	if !clientCredentialsJWT {
		userMgr := user.NewManager(r)
		var err error
		authorization, err = userMgr.GetAuthorization(username, clientID)
		if err != nil {
			log.Error("Failed to get authorization: ", err)
			return grantedScopes
		}
	}
	fullScopes = append(fullScopes, getActualValues(r, grantedScopes, username, authorization)...)
	return fullScopes
}

// getActualValues stores some acutal values for specified scopes in the jwt
// if possible. The label is parsed from the jwt, then a lookup is performed in the authorization
// backing this jwt. This is required because the jwt might have less scopes than the authorization,
// and we only include values that where in the jwt already
func getActualValues(r *http.Request, grantedScopes []string, username string, authorization *user.Authorization) []string {
	vals := []string{}
	userMgr := user.NewManager(r)
	userObj, err := userMgr.GetByName(username)
	if err != nil {
		log.Error("Failed to get user: ", err)
		return vals
	}
	// Similar to the user info API call
	valMgr := validation.NewManager(r)
	for _, scope := range grantedScopes {
		if scope == "user:name" {
			vals = append(vals, fmt.Sprintf("[%v]:%v %v", scope, userObj.Firstname, userObj.Lastname))
		}
		if strings.HasPrefix(scope, "user:email") {
			requestedLabel := strings.TrimPrefix(scope, "user:email:")
			if requestedLabel == "" || requestedLabel == "user:email" {
				requestedLabel = "main"
			}
			realLabel := getRealLabel(requestedLabel, "email", authorization)
			email, err := userObj.GetEmailAddressByLabel(realLabel)
			if err == nil {
				vals = append(vals, fmt.Sprintf("[%v]:%v", scope, email.EmailAddress))
			} else {
				log.Debug(err)
			}
		}
		if strings.HasPrefix(scope, "user:phone") {
			requestedLabel := strings.TrimPrefix(scope, "user:phone:")
			if requestedLabel == "" || requestedLabel == "user:phone" {
				requestedLabel = "main"
			}
			realLabel := getRealLabel(requestedLabel, "phone", authorization)
			phonenumber, err := userObj.GetPhonenumberByLabel(realLabel)
			if err == nil {
				vals = append(vals, fmt.Sprintf("[%v]:%v", scope, phonenumber.Phonenumber))
			} else {
				log.Debug(err)
			}
		}
		if strings.HasPrefix(scope, "user:validated:phone") {
			requestedLabel := strings.TrimPrefix(scope, "user:validated:phone:")
			if requestedLabel == "" || requestedLabel == "user:validated:phone" {
				requestedLabel = "main"
			}
			realLabel := getRealLabel(requestedLabel, "validatedphone", authorization)
			phone, err := userObj.GetPhonenumberByLabel(realLabel)
			if err == nil {
				validated, err := valMgr.IsPhonenumberValidated(authorization.Username, phone.Phonenumber)
				if err != nil {
					log.Error("Failed to verify if phone number is validated for this user: ", err)
					continue
				}
				if !validated {
					continue
				}
				vals = append(vals, fmt.Sprintf("[%v]:%v", scope, phone.Phonenumber))
			} else {
				log.Debug(err)
			}
		}
		if strings.HasPrefix(scope, "user:validated:email") {
			requestedLabel := strings.TrimPrefix(scope, "user:validated:email:")
			if requestedLabel == "" || requestedLabel == "user:validated:email" {
				requestedLabel = "main"
			}
			realLabel := getRealLabel(requestedLabel, "validatedemail", authorization)
			email, err := userObj.GetEmailAddressByLabel(realLabel)
			if err == nil {
				validated, err := valMgr.IsEmailAddressValidated(authorization.Username, email.EmailAddress)
				if err != nil {
					log.Error("Failed to verify if email address is validated for this user: ", err)
					continue
				}
				if !validated {
					continue
				}
				vals = append(vals, fmt.Sprintf("[%v]:%v", scope, email.EmailAddress))
			} else {
				log.Debug(err)
			}
		}
	}
	return vals
}

// getRealLabel loads the real label from an authorization
func getRealLabel(requestedLabel string, t string, authorization *user.Authorization) string {
	if authorization == nil {
		return requestedLabel
	}
	switch t {
	case "email":
		for _, m := range authorization.EmailAddresses {
			if requestedLabel == m.RequestedLabel {
				return m.RealLabel
			}
		}
		break
	case "phone":
		for _, m := range authorization.Phonenumbers {
			if requestedLabel == m.RequestedLabel {
				return m.RealLabel
			}
		}
		break
	case "validatedemail":
		for _, m := range authorization.ValidatedEmailAddresses {
			if requestedLabel == m.RequestedLabel {
				return m.RealLabel
			}
		}
		break
	case "validatedphone":
		for _, m := range authorization.ValidatedEmailAddresses {
			if requestedLabel == m.RequestedLabel {
				return m.RealLabel
			}
		}
		break
	}
	return ""
}

// getGrants returns a list of all the grants for a user
func getGrants(username, clientID string, r *http.Request) ([]string, error) {

	// Add grants
	grantMgr := grants.NewManager(r)
	rawGrants, err := grantMgr.GetGrantsForUser(username, clientID)
	if err != nil && !db.IsNotFound(err) {
		return nil, err
	}
	if err != nil {
		rawGrants = &grants.SavedGrants{Grants: []grants.Grant{}}
	}

	fullGrants := make([]string, len(rawGrants.Grants))
	for i, g := range rawGrants.Grants {
		fullGrants[i] = grants.FullName(g)
	}

	log.Debug("Found ", len(fullGrants), " grant(s)")
	return fullGrants, nil
}
