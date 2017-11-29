package security

import (
	"crypto/ecdsa"
	"net/http"
	"strings"
)

// OAuth2Middleware defines the common oauth2 functionality
type OAuth2Middleware struct {
	Scopes []string
}

//JWTPublicKey has the public key of the allowed JWT issuer
var JWTPublicKey *ecdsa.PublicKey

//GetAccessToken returns the access token from the authorization header or from the query parameters.
// If the authorization header starts with "bearer", "" is returned
func (om *OAuth2Middleware) GetAccessToken(r *http.Request) string {

	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		accessTokenQueryParameter := r.URL.Query().Get("access_token")
		return accessTokenQueryParameter
	}
	if strings.HasPrefix(authorizationHeader, "bearer ") || strings.HasPrefix(authorizationHeader, "Bearer ") {
		return ""
	}
	accessToken := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "token"))
	return accessToken
}
