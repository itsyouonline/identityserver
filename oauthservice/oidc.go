package oauthservice

import (
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/itsyouonline/identityserver/db/user"
)

const (
	// scopeOpenID is the mandatory scope for all OpenID Connect OAuth2 requests.
	scopeOpenID = "openid"
)

// getIDTokenFromCode returns an ID token if scopes associated with the code match
// the OpenId scope
// If no openId scope is found, the returned string is empty
// if no error, the int returned represents http.StatusOK
func getIDTokenFromCode(code string, jwtSigningKey *ecdsa.PrivateKey, r *http.Request, at *AccessToken, mgr *Manager) (string, int) {
	// get scopes
	fmt.Println(at.Scope)
	ar, err := mgr.getAuthorizationRequest(code)
	if err != nil {
		log.Debugf("something went wrong getting authorize request for the ID token: %s", err)
		return "", http.StatusInternalServerError
	}
	scopeStr := ar.Scope

	if !scopePresent(scopeStr, scopeOpenID) {
		return "", http.StatusOK
	}

	token, err := getIDTokenStr(jwtSigningKey, r, at, scopeStr)
	if err != nil {
		log.Debugf("something went wrong getting ID token: %s", err)
		return "", http.StatusBadRequest
	}

	return token, http.StatusOK
}

// scopePresent returns true if scope is in the scope string
// The scope string is expected to be a comma seperated list of scopes
func scopePresent(scopeStr string, scopeToSearch string) bool {
	scopeSlice := strings.Split(scopeStr, ",")
	for _, scope := range scopeSlice {
		if scope == scopeToSearch {
			return true
		}
	}

	return false
}

// getIDTokenStr returns an oidc ID token string
// It will set the default required claims
// and calls setValuesFromScope to set additional claims
func getIDTokenStr(jwtSigningKey *ecdsa.PrivateKey, r *http.Request, at *AccessToken, scopeStr string) (string, error) {
	// for each valid oidc standard scope, fetch related data
	token := jwt.New(jwt.SigningMethodES384)

	// setup basic claims
	token.Claims["sub"] = at.Username
	token.Claims["iss"] = issuer
	token.Claims["iat"] = at.CreatedAt.Unix()
	token.Claims["exp"] = at.ExpirationTime().Unix()
	token.Claims["aud"] = at.ClientID

	// check scopes for additional claims
	err := setValuesFromScope(token, scopeStr, r, at)
	if err != nil {
		return "", fmt.Errorf("failed to get additional claims for id token: %s", err)
	}

	return token.SignedString(jwtSigningKey)
}

// setValuesFromScope check the scopes for additional claims to be added to the provided token
func setValuesFromScope(token *jwt.Token, scopeStr string, r *http.Request, at *AccessToken) error {
	userMgr := user.NewManager(r)
	authorization, err := userMgr.GetAuthorization(at.Username, at.ClientID)
	if err != nil {
		return fmt.Errorf("failed to get authorization: %s", err)
	}
	userObj, err := userMgr.GetByName(at.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %s", err)
	}

	scopeSlice := strings.Split(scopeStr, ",")
	for _, scope := range scopeSlice {
		switch {
		case scope == "user:name":
			token.Claims[scope] = fmt.Sprintf("%s %s", userObj.Firstname, userObj.Lastname)
		case strings.HasPrefix(scope, "user:email"):
			requestedLabel := strings.TrimPrefix(scope, "user:email")
			if requestedLabel == "" || requestedLabel == "user:email" {
				requestedLabel = "main"
			}
			label := getRealLabel(requestedLabel, "email", authorization)
			email, err := userObj.GetEmailAddressByLabel(label)
			if err != nil {
				return fmt.Errorf("could not get user's email: %s", err)
			}
			token.Claims[scope] = email.EmailAddress
		case strings.HasPrefix(scope, "user:validated:email"):
			requestedLabel := strings.TrimPrefix(scope, "user:validated:email")
			if requestedLabel == "" || requestedLabel == "user:validated:email" {
				requestedLabel = "main"
			}
			label := getRealLabel(requestedLabel, "validatedemail", authorization)
			email, err := userObj.GetEmailAddressByLabel(label)
			if err != nil {
				return fmt.Errorf("could not get user's email: %s", err)
			}
			token.Claims[scope] = email.EmailAddress
		case strings.HasPrefix(scope, "user:phone"):
			requestedLabel := strings.TrimPrefix(scope, "user:phone:")
			if requestedLabel == "" || requestedLabel == "user:phone" {
				requestedLabel = "main"
			}
			label := getRealLabel(requestedLabel, "phone", authorization)
			phone, err := userObj.GetPhonenumberByLabel(label)
			if err != nil {
				return fmt.Errorf("could not get user's phone: %s", err)
			}
			token.Claims[scope] = phone.Phonenumber
		case strings.HasPrefix(scope, "user:validated:phone"):
			requestedLabel := strings.TrimPrefix(scope, "user:validated:phone:")
			if requestedLabel == "" || requestedLabel == "user:validated:phone" {
				requestedLabel = "main"
			}
			label := getRealLabel(requestedLabel, "validatedphone", authorization)
			phone, err := userObj.GetPhonenumberByLabel(label)
			if err != nil {
				return fmt.Errorf("could not get user's phone: %s", err)
			}
			token.Claims[scope] = phone.Phonenumber
		}
	}

	return nil
}
