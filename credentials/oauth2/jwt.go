package oauth2

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

//GetValidJWT returns a validated ES384 signed jwt from the authorization header that needs to start with "bearer "
// If no jwt is found in the authorization header, nil is returned
// Validation against the supplied publickey is performed
func GetValidJWT(r *http.Request, publicKey *ecdsa.PublicKey) (token *jwt.Token, err error) {
	authorizationHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorizationHeader, "bearer ") && !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return
	}
	jwtstring := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(authorizationHeader, "Bearer"), "bearer"))
	token, err = jwt.Parse(jwtstring, func(token *jwt.Token) (interface{}, error) {
		m, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		if token.Header["alg"] != m.Alg() {
			return nil, fmt.Errorf("Unexpected signing algorithm: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err == nil && !token.Valid {
		err = errors.New("Invalid jwt supplied:" + jwtstring)
	}
	return
}

func GetScopesFromJWT(token *jwt.Token) (scopes []string) {
	if token == nil {
		return
	}
	//Ignore the errors for now, we only parse valid tokens
	scopes = make([]string, 0, 0)
	rawclaims, _ := token.Claims["scope"].([]interface{})
	for _, rawclaim := range rawclaims {
		scope, _ := rawclaim.(string)
		scopes = append(scopes, scope)
	}
	return
}

// GetScopestringFromJWT turns the scopes from a jwt in to a commaseperated scopestring
func GetScopestringFromJWT(token *jwt.Token) (scopestring string) {
	if token == nil {
		return
	}
	scopes := GetScopesFromJWT(token)
	scopestring = strings.Join(scopes, ",")
	return
}

// IgnoreExpired checks if the input error is only an expired error. Nil is returned in
// this case, else the original error
func IgnoreExpired(err error) error {
	vErr, ok := err.(*jwt.ValidationError)
	if ok && vErr.Errors == jwt.ValidationErrorExpired {
		return nil
	}
	return err
}
