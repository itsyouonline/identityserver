package organization

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/credentials/oauth2"
	"github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/identityservice/security"
	"github.com/itsyouonline/identityserver/oauthservice"
)

const ItsyouonlineClientID = "itsyouonline"

// Oauth2oauth_2_0Middleware is oauth2 middleware for oauth_2_0
type Oauth2oauth_2_0Middleware struct {
	security.OAuth2Middleware
}

// newOauth2oauth_2_0Middlewarecreate new Oauth2oauth_2_0Middleware struct
func newOauth2oauth_2_0Middleware(scopes []string) *Oauth2oauth_2_0Middleware {
	om := Oauth2oauth_2_0Middleware{}
	om.Scopes = scopes

	return &om
}

// CheckScopes checks whether user has needed scopes
func (om *Oauth2oauth_2_0Middleware) CheckScopes(scopes []string) bool {
	if len(om.Scopes) == 0 {
		return true
	}

	for _, allowed := range om.Scopes {
		for _, scope := range scopes {
			if scope == allowed {
				return true
			}
		}
	}
	return false
}

func scopeStringContainsScope(scopestring, scope string) bool {
	for _, availablescope := range strings.Split(scopestring, ",") {
		availablescope = strings.Trim(availablescope, " ")
		if scope == availablescope {
			return true
		}
	}
	return false
}

// Handler return HTTP handler representation of this middleware
func (om *Oauth2oauth_2_0Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var scopes []string
		protectedOrganization := mux.Vars(r)["globalid"]
		var atscopestring string
		var username string
		var clientID string
		var globalID string

		accessToken := om.GetAccessToken(r)

		token, err := oauth2.GetValidJWT(r, security.JWTPublicKey)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if token != nil {
			// The token can either be for a user or an organization.
			// The uncaught second assertion return value ensures a missing value is translated into an empty string
			username, _ = token.Claims["username"].(string)
			globalID, _ = token.Claims["globalid"].(string)
			clientID = token.Claims["azp"].(string)
			atscopestring = oauth2.GetScopestringFromJWT(token)
		} else if accessToken != "" {
			//TODO: cache
			oauthMgr := oauthservice.NewManager(r)
			at, err := oauthMgr.GetAccessToken(accessToken)
			if err != nil {
				log.Error("Error while getting access token: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if at == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			globalID = at.GlobalID
			username = at.Username
			atscopestring = at.Scope
			clientID = at.ClientID
		} else {
			if webuser, ok := context.GetOk(r, "webuser"); ok {
				if parsedusername, ok := webuser.(string); ok && parsedusername != "" {
					username = parsedusername
					atscopestring = "admin"
					clientID = ItsyouonlineClientID
				}
			}
		}
		if (username == "" && globalID == "") || clientID == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		context.Set(r, "authenticateduser", username)
		//If the authorized organization is the protected organization itself or is a parent of it
		if len(globalID) > 0 && (globalID == protectedOrganization || strings.HasPrefix(protectedOrganization, globalID+".")) {
			scopes = []string{atscopestring}
		} else {
			orgMgr := organization.NewManager(r)
			isOwner, err := orgMgr.IsOwner(protectedOrganization, username)
			if err != nil {
				log.Error("Error while checking if user is owner of organization: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if isOwner && ((clientID == ItsyouonlineClientID && atscopestring == "admin") || scopeStringContainsScope(atscopestring, "user:admin")) {
				scopes = []string{"organization:owner"}
			} else {
				isMember, err := orgMgr.IsMember(protectedOrganization, username)
				if err != nil {
					log.Error(err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				if isMember && ((clientID == ItsyouonlineClientID && atscopestring == "admin") || scopeStringContainsScope(atscopestring, "user:admin")) {
					scopes = []string{"organization:member"}
				}
			}
		}

		//TODO: scopes "organization:info", "organization:contracts:read"

		log.Debug("Available scopes: ", scopes)

		// check scopes
		if !om.CheckScopes(scopes) {
			w.WriteHeader(403)
			return
		}

		next.ServeHTTP(w, r)
	})
}
