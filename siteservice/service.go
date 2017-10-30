package siteservice

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/itsyouonline/identityserver/communication"
	"github.com/itsyouonline/identityserver/siteservice/apiconsole"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/assets"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/components"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/html"
	"github.com/itsyouonline/identityserver/siteservice/website/packaged/thirdpartyassets"
	"github.com/itsyouonline/identityserver/specifications"
	"github.com/itsyouonline/identityserver/validation"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/credentials/totp"
	"github.com/itsyouonline/identityserver/identityservice"
	"github.com/itsyouonline/identityserver/tools/assetfs"
)

//Service is the identityserver http service
type Service struct {
	Sessions                      map[SessionType]*sessions.CookieStore
	smsService                    communication.SMSService
	phonenumberValidationService  *validation.IYOPhonenumberValidationService
	EmailService                  communication.EmailService
	emailaddressValidationService *validation.IYOEmailAddressValidationService
	version                       string
	testEnv                       bool
	identityService               *identityservice.Service
}

//NewService creates and initializes a Service
func NewService(cookieSecret string, smsService communication.SMSService, emailService communication.EmailService,
	identityservice *identityservice.Service, version string, testEnv bool) (service *Service) {
	service = &Service{smsService: smsService}

	p := &validation.IYOPhonenumberValidationService{SMSService: smsService}
	service.phonenumberValidationService = p
	e := &validation.IYOEmailAddressValidationService{EmailService: emailService}
	service.emailaddressValidationService = e

	service.identityService = identityservice

	service.version = version

	service.testEnv = testEnv

	service.initializeSessions(cookieSecret)
	return
}

//InitModels initialize persistance models
func (service *Service) InitModels() {
	service.initLoginModels()
	service.initRegistrationModels()
}

//AddRoutes registers the http routes with the router
func (service *Service) AddRoutes(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(service.HomePage)
	//Registration form
	router.Methods("GET").Path("/validateusername").HandlerFunc(service.ValidateUsername)
	router.Methods("GET").Path("/register").HandlerFunc(service.ShowRegistrationForm)
	router.Methods("POST").Path("/register").HandlerFunc(service.ProcessRegistrationForm)
	router.Methods("GET").Path("/phonevalidation").HandlerFunc(service.PhonenumberValidation)
	router.Methods("GET").Path("/pvl").HandlerFunc(service.PhonenumberValidationAndLogin)
	router.Methods("GET").Path("/emailvalidation").HandlerFunc(service.EmailValidation)
	router.Methods("POST").Path("/register/resendsms").HandlerFunc(service.ResendPhonenumberConfirmation)
	router.Methods("GET").Path("/register/smsconfirmed").HandlerFunc(service.CheckRegistrationSMSConfirmation)
	router.Methods("GET").Path("/register/emailconfirmed").HandlerFunc(service.CheckRegistrationEmailConfirmation)
	router.Methods("POST").Path("/register/smsconfirmation").HandlerFunc(service.ProcessPhonenumberConfirmationForm)
	router.Methods("POST").Path("/register/validation").HandlerFunc(service.ValidateInfo)
	router.Methods("POST").Path("/register/resendvalidation").HandlerFunc(service.ResendValidationInfo)
	//Enable us to "forget" users in case we are not in production
	router.Methods("GET").Path("/register/delete").HandlerFunc(service.ServeForgetAccountPage)
	router.Methods("POST").Path("/register/delete").HandlerFunc(service.ForgetAccountHandler)
	//Login forms
	router.Methods("GET").Path("/login").HandlerFunc(service.ShowLoginForm)
	router.Methods("POST").Path("/login").HandlerFunc(service.ProcessLoginForm)
	router.Methods("GET").Path("/login/twofamethods").HandlerFunc(service.GetTwoFactorAuthenticationMethods)
	router.Methods("POST").Path("/login/totpconfirmation").HandlerFunc(service.ProcessTOTPConfirmation)
	router.Methods("POST").Path("/login/smscode/{phoneLabel}").HandlerFunc(service.GetSmsCode)
	router.Methods("POST").Path("/login/smsconfirmation").HandlerFunc(service.Process2FASMSConfirmation)
	router.Methods("POST").Path("/login/resendsms").HandlerFunc(service.LoginResendPhonenumberConfirmation)
	router.Methods("GET").Path("/sc").HandlerFunc(service.MobileSMSConfirmation)
	router.Methods("GET").Path("/login/smsconfirmed").HandlerFunc(service.Check2FASMSConfirmation)
	router.Methods("POST").Path("/login/validateemail").HandlerFunc(service.ValidateEmail)
	router.Methods("POST").Path("/login/forgotpassword").HandlerFunc(service.ForgotPassword)
	router.Methods("POST").Path("/login/resetpassword").HandlerFunc(service.ResetPassword)
	router.Methods("GET").Path("/login/organizationinvitation/{code}").HandlerFunc(service.GetOrganizationInvitation)
	//Authorize form
	router.Methods("GET").Path("/authorize").HandlerFunc(service.ShowAuthorizeForm)
	//Facebook callback
	router.Methods("GET").Path("/facebook_callback").HandlerFunc(service.FacebookCallback)
	//Github callback
	router.Methods("GET").Path("/github_callback").HandlerFunc(service.GithubCallback)
	//Logout link
	router.Methods("GET").Path("/logout").HandlerFunc(service.Logout)
	//Error page
	router.Methods("GET").Path("/error").HandlerFunc(service.ErrorPage)
	router.Methods("GET").Path("/error{errornumber}").HandlerFunc(service.ErrorPage)
	router.Methods("GET").Path("/config").HandlerFunc(service.GetConfig)

	router.Methods("GET").Path("/version").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Version string
		}{
			Version: service.version,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&response)
	})

	router.Methods("GET").Path("/location").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the header from the cloudflare IP Geolocation service. Value is the
		// country code in ISO 3166-1 Alpha 2 format.
		location := r.Header.Get("CF-IPCountry")
		log.Debug("request location: ", location)
		response := struct {
			Location string `json:"location"`
		}{
			Location: location,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&response)
	})

	//host the assets used in the htmlpages
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(
		&assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, AssetInfo: assets.AssetInfo})))
	router.PathPrefix("/thirdpartyassets/").Handler(http.StripPrefix("/thirdpartyassets/", http.FileServer(
		&assetfs.AssetFS{Asset: thirdpartyassets.Asset, AssetDir: thirdpartyassets.AssetDir, AssetInfo: thirdpartyassets.AssetInfo})))
	router.PathPrefix("/components/").Handler(http.StripPrefix("/components/", http.FileServer(
		&assetfs.AssetFS{Asset: components.Asset, AssetDir: components.AssetDir, AssetInfo: components.AssetInfo})))

	//host the apidocumentation
	router.Methods("GET").Path("/apidocumentation").HandlerFunc(service.APIDocs)
	router.PathPrefix("/apidocumentation/raml/").Handler(http.StripPrefix("/apidocumentation/raml", http.FileServer(
		&assetfs.AssetFS{Asset: specifications.Asset, AssetDir: specifications.AssetDir, AssetInfo: specifications.AssetInfo})))
	router.PathPrefix("/apidocumentation/").Handler(http.StripPrefix("/apidocumentation/", http.FileServer(
		&assetfs.AssetFS{Asset: apiconsole.Asset, AssetDir: apiconsole.AssetDir, AssetInfo: apiconsole.AssetInfo})))

}

const (
	mainpageFileName      = "index.html"
	homepageFileName      = "base.html"
	errorpageFilename     = "error.html"
	apidocsPageFilename   = "apidocumentation.html"
	smsconfirmationPage   = "smsconfirmation.html"
	emailconfirmationPage = "emailconfirmation.html"
)

//ShowPublicSite shows the public website
func (service *Service) ShowPublicSite(w http.ResponseWriter, request *http.Request) {
	htmlData, err := html.Asset(mainpageFileName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(htmlData)
}

//APIDocs shows the api documentation
func (service *Service) APIDocs(w http.ResponseWriter, request *http.Request) {
	htmlData, err := html.Asset(apidocsPageFilename)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(htmlData)
}

//HomePage shows the home page when logged in, if not, delegate to showing the public website
func (service *Service) HomePage(w http.ResponseWriter, request *http.Request) {

	loggedinuser, err := service.GetLoggedInUser(request, w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if loggedinuser == "" {
		service.ShowPublicSite(w, request)
		return
	}

	htmlData, err := html.Asset(homepageFileName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sessions.Save(request, w)
	w.Write(htmlData)
}

//Logout logs out the user and redirect to the homepage
//TODO: csrf protection, really important here!
func (service *Service) Logout(w http.ResponseWriter, request *http.Request) {
	service.SetLoggedInUser(w, request, "")
	sessions.Save(request, w)
	http.Redirect(w, request, "", http.StatusFound)
}

//ErrorPage shows the errorpage
func (service *Service) ErrorPage(w http.ResponseWriter, request *http.Request) {
	errornumber := mux.Vars(request)["errornumber"]
	log.Debug("Errorpage requested for error ", errornumber)

	htmlData, err := html.Asset(errorpageFilename)
	if err != nil {
		log.Error("ERROR rendering error page: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// check if the error is a number to prevent XSS attacks
	errorcode, err := strconv.Atoi(errornumber)
	if err != nil {
		log.Info("Error code could not be converted to int")
		// The error page already loaded so we might as well use it
		errorcode = 400
		errornumber = "400"
	}

	// check if the error code is within the accepted 4xx client or 5xx server error range
	if errorcode > 599 || errorcode < 400 {
		log.Info("Error code out of bounds: ", errorcode)
		errorcode = 400
		errornumber = "400"
	}

	// now that we confirmed the error code is valid, we can safely use it to display on the error page
	htmlData = bytes.Replace(htmlData, []byte(`500`), []byte(errornumber), 1)
	w.Write(htmlData)
}

//renderSMSConfirmationPage renders a small mobile friendly confirmation page after a user follows a link in an sms
func (service *Service) renderSMSConfirmationPage(w http.ResponseWriter, request *http.Request, text string) {
	htmlData, err := html.Asset(smsconfirmationPage)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	htmlData = bytes.Replace(htmlData, []byte(`{{ text }}`), []byte(text), 1)
	sessions.Save(request, w)
	w.Write(htmlData)
}

//renderEmailConfirmationPage renders a small mobile friendly confirmation page after a user follows a link in an email
func (service *Service) renderEmailConfirmationPage(w http.ResponseWriter, request *http.Request, text string, extratext string) {
	htmlData, err := html.Asset(emailconfirmationPage)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	htmlData = bytes.Replace(htmlData, []byte(`{{ text }}`), []byte(text), 1)
	htmlData = bytes.Replace(htmlData, []byte(`{{ extratext }}`), []byte(extratext), 1)
	sessions.Save(request, w)
	w.Write(htmlData)
}

func (service *Service) GetConfig(w http.ResponseWriter, request *http.Request) {
	token, err := totp.NewToken()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totpsession, err := service.GetSession(request, SessionForRegistration, "totp")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totpsession.Values["secret"] = token.Secret
	sessions.Save(request, w)
	data := struct {
		TotpIssuer       string `json:"totpissuer"`
		TotpSecret       string `json:"totpsecret"`
		GithubClientId   string `json:"githubclientid"`
		FacebookClientId string `json:"facebookclientid"`
	}{
		TotpIssuer: totp.GetIssuer(request),
		TotpSecret: token.Secret,
	}
	data.GithubClientId, _ = identityservice.GetOauthClientID("github")
	data.FacebookClientId, _ = identityservice.GetOauthClientID("facebook")
	json.NewEncoder(w).Encode(&data)
}
