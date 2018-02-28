package oauthservice

import (
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/itsyouonline/identityserver/db/user"
)

const (
	// ScopeOpenID is the mandatory scope for all OpenID Connect OAuth2 requests.
	ScopeOpenID = "openid"
)

// getIDTokenFromCode returns an ID token if scopes associated with the code match
// the OpenId scope
// If no openId scope is found, the returned string is empty
func getIDTokenFromCode(code string, jwtSigningKey *ecdsa.PrivateKey, r *http.Request, at *AccessToken, mgr *Manager) (string, error) {
	// get scopes
	ar, err := mgr.getAuthorizationRequest(code)
	if err != nil {
		return "", err
	}
	scopeStr := ar.Scope

	if !findScope(scopeStr, ScopeOpenID) {
		return "", nil
	}

	return getIDToken(jwtSigningKey, r, at, scopeStr)
}

// findScope returns true if scope is in the scope string
// The scope string is expected to be a comma seperated list of scopes
func findScope(scopeStr string, scopeToSearch string) bool {
	scopeSlice := strings.Split(scopeStr, ",")
	for _, scope := range scopeSlice {
		if scope == scopeToSearch {
			return true
		}
	}

	return false
}

// getIDToken returns an oidc ID token string
// It will check the provided scopes string for supported oidc standard scope values
// and set the corresponding standard claims if available.
func getIDToken(jwtSigningKey *ecdsa.PrivateKey, r *http.Request, at *AccessToken, scopeStr string) (string, error) {
	// for each valid oidc standard scope, fetch related data
	token := jwt.New(jwt.SigningMethodES384)

	// setup basic claims
	token.Claims["sub"] = at.Username
	token.Claims["iss"] = issuer
	token.Claims["iat"] = at.CreatedAt
	token.Claims["exp"] = at.ExpirationTime().Unix()
	token.Claims["aud"] = at.ClientID

	// check scopes for additional claims
	userMgr := user.NewManager(r)
	authorization, err := userMgr.GetAuthorization(at.Username, at.ClientID)
	if err != nil {
		return "", fmt.Errorf("Failed to get authorization: %s", err)
	}
	userObj, err := userMgr.GetByName(at.Username)
	if err != nil {
		return "", fmt.Errorf("Failed to get user: %s", err)
	}

	scopeSlice := strings.Split(scopeStr, ",")
	for _, scope := range scopeSlice {
		switch scope {
		case "email":
			label := getRealLabel("main", "validatedemail", authorization)
			err := setEmailClaims(token, userObj, label)
			if err != nil {
				return "", err
			}
		case "profile":
			err := setProfileClaims(token, userObj)
			if err != nil {
				return "", err
			}
		case "phone":
			label := getRealLabel("main", "validatedphone", authorization)
			err := setPhoneClaims(token, userObj, label)
			if err != nil {
				return "", err
			}
		}
	}

	return token.SignedString(jwtSigningKey)
}

// setEmailClaims sets email claims into provided token
func setEmailClaims(token *jwt.Token, user *user.User, label string) error {
	var err error
	token.Claims["email"], err = user.GetEmailAddressByLabel(label)
	if err != nil {
		return fmt.Errorf("could not get user's email: %s", err)
	}

	token.Claims["email_verified"] = true

	return nil
}

// setPhoneClaims sets phone claims into provided token
func setPhoneClaims(token *jwt.Token, user *user.User, label string) error {
	var err error
	token.Claims["phone_number"], err = user.GetPhonenumberByLabel(label)
	if err != nil {
		return fmt.Errorf("could not get user's email: %s", err)
	}

	token.Claims["phone_number_verified"] = true

	return nil
}

// setProfileClaims sets profile claims into provided token
func setProfileClaims(token *jwt.Token, user *user.User) error {
	token.Claims["profile"] = "find fields iyo can provide here and add these as claims"
	token.Claims["given_name"] = user.Firstname
	token.Claims["family_name"] = user.Lastname
	token.Claims["name"] = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)

	return nil
}
