package routes

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/identityservice"
	"github.com/itsyouonline/identityserver/oauthservice"
	"github.com/itsyouonline/identityserver/siteservice"
)

//GetRouter contructs the router hierarchy and registers all handlers and middleware
func GetRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	cookieSecret := identityservice.GetCookieSecret()

	siteservice := siteservice.NewService(cookieSecret)
	siteservice.AddRoutes(r)

	apiRouter := r.PathPrefix("/api").Subrouter()
	identityservice.NewService().AddRoutes(apiRouter)

	oauthservice.NewService(siteservice).AddRoutes(r)

	// Add middlewares
	router := NewRouter(r)

	dbmw := db.DBMiddleware()
	recovery := handlers.RecoveryHandler()

	router.Use(recovery, LoggingMiddleware, dbmw)

	return router.Handler()
}
