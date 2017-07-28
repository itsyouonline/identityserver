package user

//This file is auto-generated by go-raml
//Do not edit this file by hand since it will be overwritten during the next generation

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// UsersInterface is interface for /users root endpoint
type UsersInterface interface { // Post is the handler for POST /users
	// Create a new user
	Post(http.ResponseWriter, *http.Request)
	// UpdateName is the handler for PUT / users/{username}/name
	UpdateName(http.ResponseWriter, *http.Request)
	// UpdatePassword is the handler for PUT /users/{username}/password
	UpdatePassword(http.ResponseWriter, *http.Request)
	// GetUserPhoneNumbers is the handler for GET /users/{username}/phonenumbers
	GetUserPhoneNumbers(http.ResponseWriter, *http.Request)
	// RegisterNewPhonenumber is the handler for POST /users/{username}/phonenumbers
	// Register a new phonenumber
	RegisterNewPhonenumber(http.ResponseWriter, *http.Request)
	// ValidatePhoneNumber is the handler for POST /users/{username}/phonenumbers/{label}/validate
	// Send sms verification to phone number
	ValidatePhoneNumber(http.ResponseWriter, *http.Request)
	// VerifyPhoneNumber is the handler for PUT /users/{username}/phonenumbers/{label}/validate
	// Verifies a phone number
	VerifyPhoneNumber(http.ResponseWriter, *http.Request)
	// GetUserPhonenumberByLabel is the handler for GET /users/{username}/phonenumbers/{label}
	GetUserPhonenumberByLabel(http.ResponseWriter, *http.Request)
	// UpdatePhonenumber is the handler for PUT /users/{username}/phonenumbers/{label}
	// Update the label and/or value of an existing phonenumber.
	UpdatePhonenumber(http.ResponseWriter, *http.Request)
	// DeletePhonenumber is the handler for DELETE /users/{username}/phonenumbers/{label}
	// Removes a phonenumber
	DeletePhonenumber(http.ResponseWriter, *http.Request)
	// GetUserBankAccounts is the handler for GET /users/{username}/banks
	GetUserBankAccounts(http.ResponseWriter, *http.Request)
	// CreateUserBankAccount is the handler for POST /users/{username}/banks
	// Create new bank account
	CreateUserBankAccount(http.ResponseWriter, *http.Request)
	// GetNotifications is the handler for GET /users/{username}/notifications
	// Get the list of notifications, these are pending invitations or approvals
	GetNotifications(http.ResponseWriter, *http.Request)
	// GetUser is the handler for GET /users/{username}
	GetUser(http.ResponseWriter, *http.Request)
	// DeleteFacebookAccount is the handler for DELETE /users/{username}/facebook
	// Delete the associated facebook account
	DeleteFacebookAccount(http.ResponseWriter, *http.Request)
	// RegisterNewEmailAddress is the handler for POST /users/{username}/emailaddresses
	// Register a new email address
	RegisterNewEmailAddress(http.ResponseWriter, *http.Request)
	// UpdateEmailAddress is the handler for PUT /users/{username}/emailaddresses/{label}
	// Updates the label and/or value of an email address
	UpdateEmailAddress(http.ResponseWriter, *http.Request)
	// DeleteEmailAddress is the handler for DELETE /users/{username}/emailaddresses/{label}
	// Removes an email address
	DeleteEmailAddress(http.ResponseWriter, *http.Request)
	ValidateEmailAddress(http.ResponseWriter, *http.Request)
	ListEmailAddresses(http.ResponseWriter, *http.Request)
	// DeleteGithubAccount is the handler for DELETE /users/{username}/github
	// Unlink Github Account
	DeleteGithubAccount(http.ResponseWriter, *http.Request)
	// GetUserInformation is the handler for GET /users/{username}/info
	GetUserInformation(http.ResponseWriter, *http.Request)
	// GetUserAddresses is the handler for GET /users/{username}/addresses
	GetUserAddresses(http.ResponseWriter, *http.Request)
	// RegisterNewAddress is the handler for POST /users/{username}/addresses
	// Register a new address
	RegisterNewAddress(http.ResponseWriter, *http.Request)
	// GetUserAddressByLabel is the handler for GET /users/{username}/addresses/{label}
	GetUserAddressByLabel(http.ResponseWriter, *http.Request)
	// UpdateAddress is the handler for PUT /users/{username}/addresses/{label}
	// Update the label and/or value of an existing address.
	UpdateAddress(http.ResponseWriter, *http.Request)
	// DeleteAddress is the handler for DELETE /users/{username}/addresses/{label}
	// Removes an address
	DeleteAddress(http.ResponseWriter, *http.Request)
	// GetUserBankAccountByLabel is the handler for GET /users/{username}/banks/{label}
	GetUserBankAccountByLabel(http.ResponseWriter, *http.Request)
	// UpdateUserBankAccount is the handler for PUT /users/{username}/banks/{label}
	// Update an existing bankaccount and label.
	UpdateUserBankAccount(http.ResponseWriter, *http.Request)
	// DeleteUserBankAccount is the handler for DELETE /users/{username}/banks/{label}
	// Delete a BankAccount
	DeleteUserBankAccount(http.ResponseWriter, *http.Request)
	// GetUserContracts is the handler for GET /users/{username}/contracts
	// Get the contracts where the user is 1 of the parties. Order descending by date.
	GetUserContracts(http.ResponseWriter, *http.Request)
	// RegisterNewContract is the handler for POST /user/{username}/contracts
	RegisterNewContract(http.ResponseWriter, *http.Request)
	// GetAllAuthorizations is the handler for GET /users/{username}/authorizations
	// Get the list of authorizations.
	GetAllAuthorizations(http.ResponseWriter, *http.Request)
	// GetAuthorization is the handler for GET /users/{username}/authorizations/{grantedTo}
	// Get the authorization for a specific organization.
	GetAuthorization(http.ResponseWriter, *http.Request)
	// UpdateAuthorization is the handler for PUT /users/{username}/authorizations/{grantedTo}
	// Modify which information an organization is able to see.
	UpdateAuthorization(http.ResponseWriter, *http.Request)
	// DeleteAuthorization is the handler for DELETE /users/{username}/authorizations/{grantedTo}
	// Remove the authorization for an organization, the granted organization will no longer
	// have access the user's information.
	DeleteAuthorization(http.ResponseWriter, *http.Request)
	// GetSeeObjects is the handler for GET /users/{username}/see
	// Get the list of see objects.
	GetSeeObjects(http.ResponseWriter, *http.Request)
	// GetSeeObjectsByOrganization is the handler for GET /users/{username}/see/{globalid}
	// Get the list of see objects for a specific organization.
	GetSeeObjectsByOrganization(http.ResponseWriter, *http.Request)
	// GetSeeObject is the handler for GET /users/{username}/see/{globalid}/{uniqueid}
	GetSeeObject(http.ResponseWriter, *http.Request)
	// CreateSeeObject is the handler for POST /users/{username}/see/{globalid}/{uniqueid}
	CreateSeeObject(http.ResponseWriter, *http.Request)
	// UpdateSeeObject is the handler for PUT /users/{username}/see/{globalid}/{uniqueid}
	UpdateSeeObject(http.ResponseWriter, *http.Request)
	// AddAPIKey Add an API Key
	AddAPIKey(http.ResponseWriter, *http.Request)
	GetAPIKey(http.ResponseWriter, *http.Request)
	UpdateAPIKey(http.ResponseWriter, *http.Request)
	DeleteAPIKey(http.ResponseWriter, *http.Request)
	ListAPIKeys(http.ResponseWriter, *http.Request)
	// AddPublicKey Add a public key
	AddPublicKey(http.ResponseWriter, *http.Request)
	// GetPublicKey Get the public key associated with a label
	GetPublicKey(http.ResponseWriter, *http.Request)
	// UpdatePublicKey Updates the label and or key of an existing public key
	UpdatePublicKey(http.ResponseWriter, *http.Request)
	// DeletePublicKey Deletes a public key
	DeletePublicKey(http.ResponseWriter, *http.Request)
	// ListPublicKeys Lists all public keys
	ListPublicKeys(http.ResponseWriter, *http.Request)
	// GetKeyStore is the handler for GET /users/{username}/keystore
	// Returns all the publickeys written to the user by an organizaton
	GetKeyStore(http.ResponseWriter, *http.Request)
	// GetKeyStoreKey is the handler for GET /users/{username}/keystore{label}
	// Returns all specific publickey written to the user by an organizaton
	GetKeyStoreKey(http.ResponseWriter, *http.Request)
	// SaveKeyStoreKey is the handler for POST /users/{username}/keystore
	// Returns all the publickeys written to the user by an organizaton
	SaveKeyStoreKey(http.ResponseWriter, *http.Request)
	// GetTwoFAMethods is the handler for GET /users/{username}/twofamethods
	// Get the possible two factor authentication methods
	GetTwoFAMethods(http.ResponseWriter, *http.Request)
	// GetTOTPSecret is the handler for GET /users/{username}/totp
	GetTOTPSecret(http.ResponseWriter, *http.Request)
	// SetupTOTP is the handler for POST /users/{username}/totp
	SetupTOTP(http.ResponseWriter, *http.Request)
	// RemoveTOTP is the handler for DELETE /users/{username}/totp
	RemoveTOTP(http.ResponseWriter, *http.Request)
	GetDigitalWallet(http.ResponseWriter, *http.Request)
	RegisterNewDigitalAssetAddress(http.ResponseWriter, *http.Request)
	GetDigitalAssetAddress(http.ResponseWriter, *http.Request)
	UpdateDigitalAssetAddress(http.ResponseWriter, *http.Request)
	DeleteDigitalAssetAddress(http.ResponseWriter, *http.Request)
	// LeaveOrganization is the handler for DELETE /users/{username}/organizations/{globalid}/leave
	LeaveOrganization(http.ResponseWriter, *http.Request)

	// ListUserRegistry is the handler for GET /users/{username}/registry
	// Lists the Registry entries
	ListUserRegistry(http.ResponseWriter, *http.Request)
	// AddUserRegistryEntry is the handler for POST /users/{username}/registry
	// Adds a RegistryEntry to the user's registry, if the key is already used, it is overwritten.
	AddUserRegistryEntry(http.ResponseWriter, *http.Request)
	// GetUserRegistryEntry is the handler for GET /users/{username}/registry/{key}
	// Get a RegistryEntry from the user's registry.
	GetUserRegistryEntry(http.ResponseWriter, *http.Request)
	// DeleteUserRegistryEntry is the handler for DELETE /users/{username}/registry/{key}
	// Removes a RegistryEntry from the user's registry
	DeleteUserRegistryEntry(http.ResponseWriter, *http.Request)
	// GetAvatar is the handler for GET /users/{username}/avatar
	// List all avatars for the user
	GetAvatars(http.ResponseWriter, *http.Request)
	// GetAvatarImage is the handler for GET /users/avatar/img/{hash}
	// Get the avatar file associated with this id
	GetAvatarImage(http.ResponseWriter, *http.Request)
	// CreateAVatarFromImage is the handler for POST /users/{username}/avatar/img/{label}
	// Create a new avatar with the specified label from a provided image file
	CreateAvatarFromImage(http.ResponseWriter, *http.Request)
	// CreateAvatarFromLink is the handler for POST /users/{username}/avatar
	// Create a new avatar with the specified label from a link
	CreateAvatarFromLink(http.ResponseWriter, *http.Request)
	// DeleteAvatar is the handler for DELETE /users/{username}/avatar/{label}
	// Delete the avatar with the specified label
	DeleteAvatar(http.ResponseWriter, *http.Request)
	// UpdateAvatarFile is the handler for PUT /users/{username}/avatar/{label}/to/{newlabel}
	// Update the avatar and possibly the avatar file stored on itsyou.online
	UpdateAvatarFile(http.ResponseWriter, *http.Request)
	// UpdateAvatarLink is the handler for PUT /users/{username}/avatar/{label}
	// Update the avatar and possibly the link to the avatar
	UpdateAvatarLink(http.ResponseWriter, *http.Request)
}

// UsersInterfaceRoutes is routing for /users root endpoint
func UsersInterfaceRoutes(r *mux.Router, i UsersInterface) {
	r.Handle("/users", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.Post))).Methods("POST")
	r.Handle("/users/{username}/name", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateName))).Methods("PUT")
	r.Handle("/users/{username}/password", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdatePassword))).Methods("PUT")
	r.Handle("/users/{username}/phonenumbers", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetUserPhoneNumbers))).Methods("GET")
	r.Handle("/users/{username}/phonenumbers", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.RegisterNewPhonenumber))).Methods("POST")
	r.Handle("/users/{username}/phonenumbers/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin", "user:phone:{label}", "user:phone:{label}:write"}).Handler).Then(http.HandlerFunc(i.GetUserPhonenumberByLabel))).Methods("GET")
	r.Handle("/users/{username}/phonenumbers/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin", "user:phone:{label}:write"}).Handler).Then(http.HandlerFunc(i.UpdatePhonenumber))).Methods("PUT")
	r.Handle("/users/{username}/phonenumbers/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin", "user:phone:{label}:write"}).Handler).Then(http.HandlerFunc(i.DeletePhonenumber))).Methods("DELETE")
	r.Handle("/users/{username}/phonenumbers/{label}/validate", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.ValidatePhoneNumber))).Methods("POST")
	r.Handle("/users/{username}/phonenumbers/{label}/validate", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.VerifyPhoneNumber))).Methods("PUT")
	r.Handle("/users/{username}/banks", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetUserBankAccounts))).Methods("GET")
	r.Handle("/users/{username}/banks", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.CreateUserBankAccount))).Methods("POST")
	r.Handle("/users/{username}/notifications", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetNotifications))).Methods("GET")
	r.Handle("/users/{username}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetUser))).Methods("GET")
	r.Handle("/users/{username}/apikeys", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.ListAPIKeys))).Methods("GET")
	r.Handle("/users/{username}/apikeys", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.AddAPIKey))).Methods("POST")
	r.Handle("/users/{username}/apikeys/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetAPIKey))).Methods("GET")
	r.Handle("/users/{username}/apikeys/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateAPIKey))).Methods("PUT")
	r.Handle("/users/{username}/apikeys/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteAPIKey))).Methods("DELETE")
	r.Handle("/users/{username}/publickeys", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.ListPublicKeys))).Methods("GET")
	r.Handle("/users/{username}/publickeys", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.AddPublicKey))).Methods("POST")
	r.Handle("/users/{username}/publickeys/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetPublicKey))).Methods("GET")
	r.Handle("/users/{username}/publickeys/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdatePublicKey))).Methods("PUT")
	r.Handle("/users/{username}/publickeys/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeletePublicKey))).Methods("DELETE")
	r.Handle("/users/{username}/keystore", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:keystore"}).Handler).Then(http.HandlerFunc(i.GetKeyStore))).Methods("GET")
	r.Handle("/users/{username}/keystore/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:keystore"}).Handler).Then(http.HandlerFunc(i.GetKeyStoreKey))).Methods("GET")
	r.Handle("/users/{username}/keystore", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:keystore"}).Handler).Then(http.HandlerFunc(i.SaveKeyStoreKey))).Methods("POST")
	r.Handle("/users/{username}/facebook", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteFacebookAccount))).Methods("DELETE")
	r.Handle("/users/{username}/emailaddresses", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.ListEmailAddresses))).Methods("GET")
	r.Handle("/users/{username}/emailaddresses", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.RegisterNewEmailAddress))).Methods("POST")
	r.Handle("/users/{username}/emailaddresses/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateEmailAddress))).Methods("PUT")
	r.Handle("/users/{username}/emailaddresses/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteEmailAddress))).Methods("DELETE")
	r.Handle("/users/{username}/emailaddresses/{label}/validate", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.ValidateEmailAddress))).Methods("POST")
	r.Handle("/users/{username}/github", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteGithubAccount))).Methods("DELETE")
	r.Handle("/users/{username}/info", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:info", "user:admin"}).Handler).Then(http.HandlerFunc(i.GetUserInformation))).Methods("GET")
	r.Handle("/users/{username}/addresses", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetUserAddresses))).Methods("GET")
	r.Handle("/users/{username}/addresses", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.RegisterNewAddress))).Methods("POST")
	r.Handle("/users/{username}/addresses/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetUserAddressByLabel))).Methods("GET")
	r.Handle("/users/{username}/addresses/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateAddress))).Methods("PUT")
	r.Handle("/users/{username}/addresses/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteAddress))).Methods("DELETE")
	r.Handle("/users/{username}/digitalwallet", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetDigitalWallet))).Methods("GET")
	r.Handle("/users/{username}/digitalwallet", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.RegisterNewDigitalAssetAddress))).Methods("POST")
	r.Handle("/users/{username}/digitalwallet/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetDigitalAssetAddress))).Methods("GET")
	r.Handle("/users/{username}/digitalwallet/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateDigitalAssetAddress))).Methods("PUT")
	r.Handle("/users/{username}/digitalwallet/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteDigitalAssetAddress))).Methods("DELETE")
	r.Handle("/users/{username}/banks/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetUserBankAccountByLabel))).Methods("GET")
	r.Handle("/users/{username}/banks/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateUserBankAccount))).Methods("PUT")
	r.Handle("/users/{username}/banks/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteUserBankAccount))).Methods("DELETE")
	r.Handle("/users/{username}/contracts", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetUserContracts))).Methods("GET")
	r.Handle("/users/{username}/contracts", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.RegisterNewContract))).Methods("POST")
	r.Handle("/users/{username}/authorizations", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetAllAuthorizations))).Methods("GET")
	r.Handle("/users/{username}/authorizations/{grantedTo}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetAuthorization))).Methods("GET")
	r.Handle("/users/{username}/authorizations/{grantedTo}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateAuthorization))).Methods("PUT")
	r.Handle("/users/{username}/authorizations/{grantedTo}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteAuthorization))).Methods("DELETE")
	r.Handle("/users/{username}/see", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetSeeObjects))).Methods("GET")
	r.Handle("/users/{username}/see/{globalid}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin", "user:see"}).Handler).Then(http.HandlerFunc(i.GetSeeObjectsByOrganization))).Methods("GET")
	r.Handle("/users/{username}/see/{globalid}/{uniqueid}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin", "user:see"}).Handler).Then(http.HandlerFunc(i.GetSeeObject))).Methods("GET")
	r.Handle("/users/{username}/see/{globalid}/{uniqueid}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:see"}).Handler).Then(http.HandlerFunc(i.CreateSeeObject))).Methods("POST")
	r.Handle("/users/{username}/see/{globalid}/{uniqueid}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:see"}).Handler).Then(http.HandlerFunc(i.UpdateSeeObject))).Methods("PUT")
	r.Handle("/users/{username}/twofamethods", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetTwoFAMethods))).Methods("GET")
	r.Handle("/users/{username}/totp", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetTOTPSecret))).Methods("GET")
	r.Handle("/users/{username}/totp", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.SetupTOTP))).Methods("POST")
	r.Handle("/users/{username}/totp", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.RemoveTOTP))).Methods("DELETE")
	r.Handle("/users/{username}/organizations/{globalid}/leave", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.LeaveOrganization))).Methods("DELETE")
	r.Handle("/users/{username}/registry", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.ListUserRegistry))).Methods("GET")
	r.Handle("/users/{username}/registry", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.AddUserRegistryEntry))).Methods("POST")
	r.Handle("/users/{username}/registry/{key}", alice.New(newUserIndentifierMiddleware().Handler).Then(http.HandlerFunc(i.GetUserRegistryEntry))).Methods("GET")
	r.Handle("/users/{username}/registry/{key}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteUserRegistryEntry))).Methods("DELETE")
	r.Handle("/users/{username}/avatar", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.GetAvatars))).Methods("GET")
	r.Handle("/users/avatar/img/{hash}", http.HandlerFunc(i.GetAvatarImage)).Methods("GET")
	r.Handle("/users/{username}/avatar/img/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.CreateAvatarFromImage))).Methods("POST")
	r.Handle("/users/{username}/avatar", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.CreateAvatarFromLink))).Methods("POST")
	r.Handle("/users/{username}/avatar/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.DeleteAvatar))).Methods("DELETE")
	r.Handle("/users/{username}/avatar/{label}/to/{newlabel}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateAvatarFile))).Methods("PUT")
	r.Handle("/users/{username}/avatar/{label}", alice.New(newUserIndentifierMiddleware().Handler, newOauth2oauth_2_0Middleware([]string{"user:admin"}).Handler).Then(http.HandlerFunc(i.UpdateAvatarLink))).Methods("PUT")

}
