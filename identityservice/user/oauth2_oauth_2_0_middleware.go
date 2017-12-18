package user

import (
	"net/http"
	"strings"

	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/credentials/oauth2"
	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/identityservice/security"
	"github.com/itsyouonline/identityserver/oauthservice"
)

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

// Handler return HTTP handler representation of this middleware
func (om *Oauth2oauth_2_0Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var atscopestring string
		var username string
		var clientID string
		authorizedScopes := []string{}

		accessToken := om.GetAccessToken(r)

		token, err := oauth2.GetValidJWT(r, security.JWTPublicKey)
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if token != nil {
			username = token.Claims["username"].(string)
			clientID = token.Claims["azp"].(string)
			atscopestring = oauth2.GetScopestringFromJWT(token)

		} else if accessToken != "" {
			//TODO: cache
			oauthMgr := oauthservice.NewManager(r)
			at, err := oauthMgr.GetAccessToken(accessToken)
			if err != nil {
				log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if at == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			username = at.Username
			atscopestring = at.Scope
			clientID = at.ClientID
		} else {
			if webuser, ok := context.GetOk(r, "webuser"); ok {
				if parsedusername, ok := webuser.(string); ok && parsedusername != "" {
					username = parsedusername
					atscopestring = "admin"
					clientID = "itsyouonline"
				}
			}
		}
		log.Debugln("Accessing user:", username, "- Accessing clientID:", clientID)
		if username == "" || clientID == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		protectedUsername := mux.Vars(r)["username"]

		for _, scope := range strings.Split(atscopestring, ",") {
			scope = strings.Trim(scope, " ")
			authorizedScopes = append(authorizedScopes, scope)
		}

		possibleScopes := []string{}
		// Replace {variables} in the scopes with the real vars present in the url.
		for _, scope := range om.Scopes {
			regex, _ := regexp.Compile(`\{(.*?)\}`)
			vars := regex.FindStringSubmatch(scope)
			if vars != nil {
				// e.g. vars -> ["{label}", "label"]
				re, _ := regexp.Compile(vars[0])
				urlValue := mux.Vars(r)[vars[1]]
				scope = re.ReplaceAllString(scope, urlValue)
			}
			possibleScopes = append(possibleScopes, scope)
		}

		// atscopestring will be user:admin for user api keys, which is only valid if the api key is owned by the user being accessed off course
		if !((protectedUsername == username && atscopestring == "user:admin") || (clientID == "itsyouonline" && atscopestring == "admin")) {
			// todo: cache
			userMgr := user.NewManager(r)
			authorization, err := userMgr.GetAuthorization(protectedUsername, clientID)
			if err != nil {
				log.Error("Error while getting authorization: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if authorization == nil {
				log.Debugf("Authorization for client id %s not found for user %s", clientID, protectedUsername)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			authorizedScopes = authorization.FilterAuthorizedScopes(authorizedScopes)
		}

		if protectedUsername == username && clientID == "itsyouonline" && atscopestring == "admin" {
			authorizedScopes = append(authorizedScopes, "user:admin")
		}
		if strings.HasPrefix(atscopestring, "user:") {
			authorizedScopes = append(authorizedScopes, "user:info")
		}

		idAzp := context.Get(r, "iyoid_azp")
		if idAzp != "" && idAzp != clientID {
			log.Debugf("Iyo id azp from the user identifier and client ID mismatch: client: %s, iyo id belongs to %s", clientID, idAzp)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		context.Set(r, "client_id", clientID)
		context.Set(r, "availablescopes", strings.Join(authorizedScopes, ","))

		// check scopes
		log.Debug("Authorized scopes: ", authorizedScopes)
		log.Debug("Needed possible scopes: ", possibleScopes)
		if !oauth2.CheckScopes(possibleScopes, authorizedScopes) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
