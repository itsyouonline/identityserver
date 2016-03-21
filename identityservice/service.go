package identityservice

import (
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/identityservice/company"
	"github.com/itsyouonline/identityserver/identityservice/globalconfig"
	"github.com/itsyouonline/identityserver/identityservice/organization"
	"github.com/itsyouonline/identityserver/identityservice/user"
	"github.com/itsyouonline/identityserver/identityservice/userorganization"

	"crypto/rand"
	"encoding/base64"

	log "github.com/Sirupsen/logrus"
)

//Service is the identityserver http service
type Service struct {
}

//NewService creates and initializes a Service
func NewService() *Service {
	return &Service{}
}

//AddRoutes registers the http routes with the router.
func (service *Service) AddRoutes(router *mux.Router) {
	// User API
	user.UsersInterfaceRoutes(router, user.UsersAPI{})
	user.InitModels()

	// Company API
	company.CompaniesInterfaceRoutes(router, company.CompaniesAPI{})
	company.InitModels()

	// Organization API
	organization.OrganizationsInterfaceRoutes(router, organization.OrganizationsAPI{})
	userorganization.UserorganizationsInterfaceRoutes(router, userorganization.UsersusernameorganizationsAPI{})
	organization.InitModels()

}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Generate a random string (s length) used for secret cookie
func generateCookieSecret(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Get secret cookie from mongodb if exists otherwise, generate a new one and save it
func GetCookieSecret() string {
	session := db.GetSession()
	defer session.Close()

	config := globalconfig.NewManager()
	globalconfig.InitModels()

	cookie, err := config.GetByKey("cookieSecret")
	if err != nil {
		log.Debug("No cookie secret found, generating a new one")

		secret, err := generateCookieSecret(32)

		if err != nil {
			log.Panic("Cannot generate secret cookie")
		}

		cookie.Key = "cookieSecret"
		cookie.Value = secret

		err = config.Insert(cookie)

		// Key was inserted by another instance in the meantime
		if db.IsDup(err) {
			cookie, err = config.GetByKey("cookieSecret")

			if err != nil {
				log.Panic("Cannot retreive cookie secret")
			}
		}
	}

	log.Debug("Cookie secret: ", cookie.Value)

	return cookie.Value
}
