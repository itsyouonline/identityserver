package user

import (
	"net/http"
	"strings"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/iyoid"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/db/validation"
)

// UserIdentifierMiddleware is the representation of a userIdentifierMiddleware
type UserIdentifierMiddleware struct{}

// NewUserIdentifierMiddleware creates a new userIdentifierMiddleware struct
func NewUserIdentifierMiddleware() *UserIdentifierMiddleware {
	return &UserIdentifierMiddleware{}
}

// Handler return HTTP handler representation of this middleware
// replaces the useridentifier, in the {username} section of the url with the
// associated username, if any
func (uim *UserIdentifierMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["username"]
		if strings.HasPrefix(username, "+") { //its a phone number
			valMgr := validation.NewManager(r)
			validatedPhoneNumber, err := valMgr.GetByPhoneNumber(username)
			if err != nil && !valMgr.IsErrNotFound(err) {
				log.Error("Failed to get validated phone number: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if err == nil {
				username = validatedPhoneNumber.Username
			}
		} else if strings.Contains(username, "@") { // its an email
			valMgr := validation.NewManager(r)
			validatedEmailAddress, err := valMgr.GetByEmailAddress(username)
			if err != nil && !valMgr.IsErrNotFound(err) {
				log.Error("Failed to get validated email address: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if err == nil {
				username = validatedEmailAddress.Username
			}
		} else { // check iyo ids
			mgr := iyoid.NewManager(r)
			idObj, err := mgr.GetByID(username)
			if err != nil && !db.IsNotFound(err) {
				log.Error("Failed to get iyo id: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if err == nil {
				username = idObj.Username
				context.Set(r, "iyoid_azp", idObj.Azp)
			}
		}
		// replace verified phone numbers and email addresses by the associated username
		mux.Vars(r)["username"] = username

		next.ServeHTTP(w, r)
	})
}
