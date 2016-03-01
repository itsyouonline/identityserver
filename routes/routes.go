package routes

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/identityservice"
	"github.com/itsyouonline/identityserver/oauthservice"
	"github.com/itsyouonline/identityserver/sessions"
	"github.com/itsyouonline/identityserver/siteservice"
)

var (
	sessionAuthKey  = []byte("TMJG63JWE28UQ6DRSPLCN7PZFWTM6B3PDFSBPZL2DV9WEKR83FMDVLR6TDK4TZIP")
	sessionBlockKey = []byte("G90QWGJU0AN6O5L15DVVRYVSUS7QTDO2")
)

func GetRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	siteservice := siteservice.NewService()
	siteservice.AddRoutes(r)
	identityservice.NewService().AddRoutes(r)
	oauthservice.NewService(siteservice).AddRoutes(r)

	// Add middlewares
	router := NewRouter(r)

	dbmw := db.DBMiddleware()
	sessionmw := sessions.SessionMiddleware(sessionAuthKey, sessionBlockKey)
	recovery := handlers.RecoveryHandler()

	router.Use(recovery, LoggingMiddleware, dbmw, sessionmw)

	return router.Handler()
}
