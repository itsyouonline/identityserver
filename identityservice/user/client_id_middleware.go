package user

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/itsyouonline/identityserver/credentials/oauth2"
	"github.com/itsyouonline/identityserver/identityservice/security"
	"github.com/itsyouonline/identityserver/oauthservice"
)

// ClientIDMiddleware is oauth2 middleware that sets the callers client_id on the request
type ClientIDMiddleware struct {
	security.OAuth2Middleware
}

// newClientIDMiddlware new ClientIDMiddlware struct
func newClientIDMiddleware(scopes []string) *ClientIDMiddleware {
	om := ClientIDMiddleware{}
	om.Scopes = scopes
	return &om
}

// Handler return HTTP handler representation of this middleware
func (om *ClientIDMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientID string

		accessToken := om.GetAccessToken(r)

		token, err := oauth2.GetValidJWT(r, security.JWTPublicKey)
		if err != nil {
			log.Error("Failed to get valid JWT: ", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if token != nil {
			clientID = token.Claims["azp"].(string)

		} else if accessToken != "" {
			oauthMgr := oauthservice.NewManager(r)
			at, err := oauthMgr.GetAccessToken(accessToken)
			if err != nil {
				log.Error("Failed to get access token: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if at == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			clientID = at.ClientID
		} else {
			if webuser, ok := context.GetOk(r, "webuser"); ok {
				if parsedusername, ok := webuser.(string); ok && parsedusername != "" {
					clientID = "itsyouonline"
				}
			}
		}
		if clientID == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		context.Set(r, "client_id", clientID)

		next.ServeHTTP(w, r)
	})
}
