package sessions

import (
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
)

func SetupRequestSession(w http.ResponseWriter, r *http.Request) {
	var session *sessions.Session

	// First, check if we have a sessionId set in the cookie.
	sessionId, err := getSessionId(r)
	if err != nil {
		// We need a new session
		CreateSession(w, r)
		return
	}

	// We have an existing session, we need to prepare the request context
	session, err = sessionStore.Get(r, sessionId)
	if err != nil {
		log.Error("Failed to get a session")
		panic(err)
	}

	setRequestSession(r, session)
}

// CreateSession create new Session and set its cookie.
func CreateSession(w http.ResponseWriter, r *http.Request) string {
	sessionId := uuid.NewV4().String()

	session, err := sessionStore.Get(r, sessionId)
	if err != nil {
		log.Error("Failed to get a session")
		panic(err)
	}

	// Set cookie in response!
	setCookie(w, sessionCookieName, sessionId)

	// Update request context with new session!
	setRequestSession(r, session)

	return sessionId
}

// RefreshSession destroy existing session and refresh its corresponding cookie.
// return newly created sessionId.
func RefreshSession(w http.ResponseWriter, r *http.Request, sessionId string) string {
	session, err := sessionStore.Get(r, sessionId)
	if err == nil {
		log.Error("Failed to get a session")
		panic(err)
	}

	// Delete existing session!
	session.Options.MaxAge = -1
	err = sessionStore.Save(r, w, session)
	if err != nil {
		log.Error("Failed to delete session")
	}

	// Create new session and return ID.
	return CreateSession(w, r)
}

//GetRequestSession return session from context.
func GetRequestSession(r *http.Request) *sessions.Session {
	return context.Get(r, sessionName).(*sessions.Session)
}

func setRequestSession(r *http.Request, session *sessions.Session) {
	context.Set(r, RequestSession, session)
}

// getSessionId return sessionId if it exists in the request cookie.
func getSessionId(r *http.Request) (string, error) {
	var sessionId string

	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		if err := cookieStore.Decode(sessionCookieName, cookie.Value, &sessionId); err == nil {
			return sessionId, nil
		} else {
			// Failed to decode cookie.
			// This might be caused by outdated keys or some replay attack!
			log.Error("Failed to decode cookie:", err.Error())
		}
	}

	return sessionId, errors.New("Cannot find cookie!")
}

func setCookie(w http.ResponseWriter, name string, value string) {
	if encoded, err := cookieStore.Encode(name, value); err == nil {
		cookie := &http.Cookie{
			Name:     name,
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
			// Secure:   true,
		}
		http.SetCookie(w, cookie)
	} else {
		panic(errors.New("Failed to set cookie!"))
	}
}
