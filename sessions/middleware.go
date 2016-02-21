// Copyright 2016 the ItsYou.online developers

package sessions

import (
	"errors"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"

	"github.com/itsyouonline/identityserver/db"
)

const (
	sessionStoreCollection = "itsyouonline-sessions"
	sessionMaxAge          = 300
	sessionTTL             = true

	sessionName       = "itsyouonline"
	sessionCookieName = "sessionid"

	RequestSession = "itsyouonline/identityserver/session"
)

type Handler struct {
	handler http.Handler
}

var sessionStore *MongoStore

var cookieStore = securecookie.New(secureCookieHashKey, secureCookieBlockKey)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if sessionStore == nil {
		panic(errors.New("Failed to create sessions. Session store is not available!"))
	}

	// Get sessionId, and set session cookie if not set!
	SetupRequestSession(w, r)

	h.handler.ServeHTTP(w, r)
}

func SessionMiddleware(keyPairs ...[]byte) func(h http.Handler) http.Handler {

	// Create a new session store once DB connection is ready.
	go createSessionStore(keyPairs...)

	return func(h http.Handler) http.Handler {
		return &Handler{
			handler: h,
		}
	}
}

//createSessionStore ensure session store is initialized in the DB.
func createSessionStore(keyPairs ...[]byte) {
	if sessionStore != nil {
		return
	}

	for {
		mgoSession := db.NewSession()
		if mgoSession != nil {
			// Get our session store collection in DB.
			coll := db.GetCollection(mgoSession, sessionStoreCollection)

			// Create the store.
			sessionStore = NewMongoStore(coll, sessionMaxAge, sessionTTL, keyPairs...)
			if sessionStore != nil {
				log.Debug("Session store initialized successfully!")
				break
			}
		}

		time.Sleep(1 * time.Second)
	}
}
