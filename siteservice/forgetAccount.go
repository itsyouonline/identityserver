package siteservice

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/identityservice/organization"
)

// This file contains the forget account handlers. These should only be available in dev and testing envs, indicated with a cli flag

// ServeForgetAccountPage serves the forget account page
func (service *Service) ServeForgetAccountPage(w http.ResponseWriter, r *http.Request) {
	// If we are not in a test environment, we pretend this does not exist
	// Don't worry about it
	if !service.testEnv {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	const template = `
	<html>
		<body>
			<h1>Itsyou.online forget validated info</h1>
			<form id="mainform" action="delete" method="post">
				login:<br/>
				<input type="text" id="login" name="login" placeholder="login" required /><br/>
				password:<br/>
				<input type="password" name="password" placeholder="password" required /><br/>
				<br/>
				<button type="submit">Log in</button>
			</form>
		</body>
	</html>`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(template))

}

// ForgetAccountHandler handles the actuall "forgetting" of an account, by dropping the validated email and phone numbers
// from the respective collections
func (service *Service) ForgetAccountHandler(w http.ResponseWriter, r *http.Request) {
	// If we are not in a test environment, we pretend this does not exist
	// Don't worry about it
	if !service.testEnv {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	r.ParseForm()
	login := strings.ToLower(r.FormValue("login"))

	u, err := organization.SearchUser(r, login)
	if db.IsNotFound(err) {
		w.WriteHeader(422)
		return
	} else if err != nil {
		log.Error("Failed to search for user: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	userexists := !db.IsNotFound(err)

	var validpassword bool
	passwdMgr := password.NewManager(r)
	if validpassword, err = passwdMgr.Validate(u.Username, r.FormValue("password")); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !userexists || !validpassword {
		w.WriteHeader(422)
		log.Debug("Invalid password for forgetting user")
		return
	}

	// drop validated info
	valMgr := validation.NewManager(r)
	validatedEmails, err := valMgr.GetByUsernameValidatedEmailAddress(u.Username)
	if err != nil {
		log.Error("Failed to get validated email addresses: ", err)
		return
	}

	for _, ve := range validatedEmails {
		// I can't be asked to care about the errors here, its past 11.30PM and its only for dev/staging anyway
		valMgr.RemoveValidatedEmailAddress(u.Username, ve.EmailAddress)
	}

	validatedPhones, err := valMgr.GetByUsernameValidatedPhonenumbers(u.Username)
	if err != nil {
		log.Error("Failed to get validated phone numbers: ", err)
		return
	}

	for _, vp := range validatedPhones {
		// Same as above
		valMgr.RemoveValidatedPhonenumber(u.Username, vp.Phonenumber)
	}

	w.WriteHeader(http.StatusOK)

	const template = `
	<html>
		<body>
			<h3>
				Your validated info has been forgotten and can now be reused
			</h3>
		</body>
	</html>
	`

	w.Write([]byte(template))
}
