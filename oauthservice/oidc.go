package oauthservice

import (
	"crypto/ecdsa"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/itsyouonline/identityserver/db/user"
)

const (
	// ScopeOpenID is the mandatory scope for all OpenID Connect OAuth2 requests.
	ScopeOpenID = "openid"
)

// isOIDC returns true if an access token should be returned as OIDC
// The scope string is expected to be a comma seperated list of scopes
func isOIDC(scopes string) bool {
	return findScope(scopes, ScopeOpenID)
}

// findScope returns true if scope is in the scope string
// The scope string is expected to be a comma seperated list of scopes
func findScope(scopes string, scopeToSearch string) bool {
	scopeSlice := strings.Split(scopes, ",")
	for _, scope := range scopeSlice {
		if scope == scopeToSearch {
			return true
		}
	}

	return false
}

// getIDToken returns an oidc ID token/claims
// It will check the provided scopes string for oidc standard scope values
// and set the corresponding standard claims if available.
func getIDToken(jwtSigningKey *ecdsa.PrivateKey, r *http.Request, at *AccessToken, scopes string, audiences string) (string, error) {
	// for each valid oidc standard scope, fetch related data
	token := jwt.New(jwt.SigningMethodES384)

	// setup basic claims
	token.Claims["sub"] = at.Username
	token.Claims["iss"] = issuer
	token.Claims["iat"] = at.CreatedAt
	token.Claims["exp"] = at.ExpirationTime().Unix()

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

	// check scopes for additional claims
	scopeSlice := strings.Split(scopes, ",")
	for _, scope := range scopeSlice {
		switch scope {
		case "email":
			setEmailClaims(token, r, at.Username)
		case "profile":
			setProfileClaims(token)
		case "phone":
			setPhoneClaims(token)
		}
	}

	return token.SignedString(jwtSigningKey)
}

// setEmailClaims sets email claims into provided token
func setEmailClaims(token *jwt.Token, r *http.Request, username string) error {
	userMgr := user.NewManager(r)
	userObj, err := userMgr.GetByName(username)
	if err != nil {
		return err
	}

	_ = userObj

	token.Claims["email"] = "fetch email"
	token.Claims["email_verified"] = true

	return nil
}

func setPhoneClaims(token *jwt.Token) error {
	token.Claims["phone_number"] = "fetch phone_number"
	token.Claims["phone_number_verified"] = true

	return nil
}

// setProfileClaims sets profile claims into provided token
func setProfileClaims(token *jwt.Token) error {
	token.Claims["profile"] = "find fields iyo can provide here and add these as claims"
	token.Claims["given_name"] = "firstname"
	token.Claims["family_name"] = "lastname"
	token.Claims["name"] = "firstname + lastname"

	return nil
}
