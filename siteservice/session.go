package siteservice

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
)

//SessionType is used to define the type of session
type SessionType int

const (
	//SessionForRegistration is the short anynymous session used during registration
	SessionForRegistration SessionType = iota
	//SessionInteractive is the session of an authenticated user on the itsyou.online website
	SessionInteractive SessionType = iota
)

func (service *Service) initializeSessions(cookieSecret string) {
	service.Sessions = make(map[SessionType]*sessions.CookieStore)

	cookieStoreSecret := cookieSecret

	registrationSessionStore := sessions.NewCookieStore([]byte(cookieStoreSecret))
	registrationSessionStore.Options.HttpOnly = true
	//TODO: enable this when we have automatic https
	//registrationSessionStore.Options.Secure = true
	registrationSessionStore.Options.MaxAge = 10 * 60 //10 minutes

	service.Sessions[SessionForRegistration] = registrationSessionStore

	interactiveSessionStore := sessions.NewCookieStore([]byte(cookieStoreSecret))
	interactiveSessionStore.Options.HttpOnly = true
	//TODO: enable this when we have automatic https
	//registrationSessionStore.Options.Secure = true
	interactiveSessionStore.Options.MaxAge = 10 * 60 //10 minutes

	service.Sessions[SessionInteractive] = registrationSessionStore

}

//GetSession returns the a session of the specified kind and a spefic name
func (service *Service) GetSession(request *http.Request, kind SessionType, name string) (*sessions.Session, error) {
	return service.Sessions[kind].Get(request, name)
}

//SetLoggedInUser creates a session for an authenticated user
func (service *Service) SetLoggedInUser(w http.ResponseWriter, request *http.Request, username string) (err error) {
	authenticatedSession, err := service.GetSession(request, SessionInteractive, "authenticatedsession")
	if err != nil {
		log.Error(err)
		return
	}
	authenticatedSession.Values["username"] = username

	// Set user cookie after successful login
	cookie := &http.Cookie{
		Name:  "itsyou.online.user",
		Path:  "/",
		Value: username,
	}
	http.SetCookie(w, cookie)

	return
}

//GetLoggedInUser returns an authenticated user, or an empty string if there is none
func (service *Service) GetLoggedInUser(request *http.Request) (username string, err error) {
	authenticatedSession, err := service.GetSession(request, SessionInteractive, "authenticatedsession")
	if err != nil {
		log.Error(err)
		return
	}
	savedusername := authenticatedSession.Values["username"]
	if savedusername != nil {
		username, _ = savedusername.(string)
	}
	return
}
