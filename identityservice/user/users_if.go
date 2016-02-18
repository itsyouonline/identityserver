package user

//This file is auto-generated by go-raml
//Do not edit this file by hand since it will be overwritten during the next generation

import (
	"github.com/gorilla/mux"
	"net/http"
)

type UsersInterface interface {

	// Create a new user
	// It is handler for POST /users
	Post(http.ResponseWriter, *http.Request)

	// It is handler for GET /users/{username}/validate
	usernamevalidateGet(http.ResponseWriter, *http.Request)

	// Update existing user. Updating ``username`` is not allowed.
	// It is handler for PUT /users/{username}
	usernamePut(http.ResponseWriter, *http.Request)

	// It is handler for GET /users/{username}/info
	usernameinfoGet(http.ResponseWriter, *http.Request)
}

func UsersInterfaceRoutes(r *mux.Router, i UsersInterface) {

	r.HandleFunc("/users", i.Post).Methods("POST")

	r.HandleFunc("/users/{username}/validate", i.usernamevalidateGet).Methods("GET")

	r.HandleFunc("/users/{username}", i.usernamePut).Methods("PUT")

	r.HandleFunc("/users/{username}/info", i.usernameinfoGet).Methods("GET")

}
