package organization

//This file is auto-generated by go-raml
//Do not edit this file by hand since it will be overwritten during the next generation

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// OrganizationsInterface is interface for /organizations root endpoint
type OrganizationsInterface interface { // CreateNewOrganization is the handler for POST /organizations
	// Create a new organization. 1 user should be in the owners list. Validation is performed
	// to check if the securityScheme allows management on this user.
	CreateNewOrganization(http.ResponseWriter, *http.Request)
	// GetOrganization is the handler for GET /organizations/{globalid}
	// Get organization info
	GetOrganization(http.ResponseWriter, *http.Request)
	// CreateNewSubOrganization is the handler for POST /organizations/{globalid}/suborganizations
	// Create a new suborganization.
	CreateNewSubOrganization(http.ResponseWriter, *http.Request)
	// UpdateOrganization is the handler for PUT /organizations/{globalid}
	// Update organization info
	UpdateOrganization(http.ResponseWriter, *http.Request)
	// DeleteOrganization is the handler for DELETE /organizations/{globalid}
	// Removes an organization and all associated data.
	DeleteOrganization(http.ResponseWriter, *http.Request)
	// GetAPIKeyLabels is the handler for GET /organizations/{globalid}/apikeys
	// Get the list of active api keys. The secrets themselves are not included.
	GetAPIKeyLabels(http.ResponseWriter, *http.Request)
	// CreateNewAPIKey is the handler for POST /organizations/{globalid}/apikeys
	// Create a new API Key, a secret itself should not be provided, it will be generated
	// serverside.
	CreateNewAPIKey(http.ResponseWriter, *http.Request)
	// GetAPIKey is the handler for GET /organizations/{globalid}/apikeys/{label}
	GetAPIKey(http.ResponseWriter, *http.Request)
	// UpdateAPIKey is the handler for PUT /organizations/{globalid}/apikeys/{label}
	// Updates the label or other properties of a key.
	UpdateAPIKey(http.ResponseWriter, *http.Request)
	// DeleteAPIKey is the handler for DELETE /organizations/{globalid}/apikeys/{label}
	// Removes an API key
	DeleteAPIKey(http.ResponseWriter, *http.Request)
	// GetOrganizationTree is the handler for GET /organizations/{globalid}/tree
	GetOrganizationTree(http.ResponseWriter, *http.Request)
	// UpdateOrganizationMemberShip is the handler for PUT /organizations/{globalid}/members
	UpdateOrganizationMemberShip(http.ResponseWriter, *http.Request)
	// AddOrganizationMember is the handler for POST /organizations/{globalid}/members
	// Assign a member to organization.
	AddOrganizationMember(http.ResponseWriter, *http.Request)
	// RemoveOrganizationMember is the handler for DELETE /organizations/{globalid}/members/{username}
	// Remove a member from organization
	RemoveOrganizationMember(http.ResponseWriter, *http.Request)
	// AddOrganizationOwner is the handler for POST /organizations/{globalid}/owners
	// Invite a user to become owner of an organization.
	AddOrganizationOwner(http.ResponseWriter, *http.Request)
	// RemoveOrganizationOwner is the handler for DELETE /organizations/{globalid}/owners/{username}
	// Remove an owner from organization
	RemoveOrganizationOwner(http.ResponseWriter, *http.Request)
	// GetContracts is the handler for GET /organizations/{globalid}/contracts
	// Get the contracts where the organization is 1 of the parties. Order descending by
	// date.
	GetContracts(http.ResponseWriter, *http.Request)
	// GetPendingInvitations is the handler for GET /organizations/{globalid}/invitations
	// Get the list of pending invitations for users to join this organization.
	GetPendingInvitations(http.ResponseWriter, *http.Request)
	// RemovePendingInvitation is the handler for DELETE /organizations/{globalid}/invitations/{username}
	// Cancel a pending invitation.
	RemovePendingInvitation(http.ResponseWriter, *http.Request)
	// CreateDns is the handler for POST /organizations/{globalid}/dns
	// Creates a new DNS name associated with an organization
	CreateDns(http.ResponseWriter, *http.Request)
	// UpdateDns is the handler for PUT /organizations/{globalid}/dns/{dnsname}
	// Updates an existing DNS name associated with an organization
	UpdateDns(http.ResponseWriter, *http.Request)
	// DeleteDNS is the handler for DELETE /organizations/{globalid}/dns/{dnsname}
	// Removes a DNS name
	DeleteDns(http.ResponseWriter, *http.Request)

	// ListOrganizationRegistry is the handler for GET /organizations/{globalid}/registry
	// Lists the Registry entries
	ListOrganizationRegistry(http.ResponseWriter, *http.Request)
	// AddOrganizationRegistryEntry is the handler for POST /organizations/{globalid}/registry
	// Adds a RegistryEntry to the organization's registry, if the key is already used, it is overwritten.
	AddOrganizationRegistryEntry(http.ResponseWriter, *http.Request)
	// GetOrganizationRegistryEntry is the handler for GET /organizations/{username}/globalid/{key}
	// Get a RegistryEntry from the organization's registry.
	GetOrganizationRegistryEntry(http.ResponseWriter, *http.Request)
	// DeleteOrganizationRegistryEntry is the handler for DELETE /organizations/{username}/globalid/{key}
	// Removes a RegistryEntry from the organization's registry
	DeleteOrganizationRegistryEntry(http.ResponseWriter, *http.Request)
	// SetOrganizationLogo is the handle for PUT /organizations/globalid/logo
	// Set the organization Logo for the organization
	SetOrganizationLogo(http.ResponseWriter, *http.Request)
	// GetOrganizationLogo is the handler for GET /organizations/globalid/logo
	// Get the Logo from an organization
	GetOrganizationLogo(http.ResponseWriter, *http.Request)
	// DeleteOrganizationLogo is the handler for DELETE /organizations/globalid/logo
	// Removes the Logo from an organization
	DeleteOrganizationLogo(http.ResponseWriter, *http.Request)
	// Get2faValidityTime is the handler for GET /organizations/globalid/2fa/validity
	// Get the 2fa validity time for the organization, in seconds
	Get2faValidityTime(w http.ResponseWriter, r *http.Request)
	// Set2faValidityTime is the handler for PUT /organizations/globalid/2fa/validity
	// Sets the 2fa validity time for the organization, in seconds
	Set2faValidityTime(w http.ResponseWriter, r *http.Request)
}

// OrganizationsInterfaceRoutes is routing for /organizations root endpoint
func OrganizationsInterfaceRoutes(r *mux.Router, i OrganizationsInterface) {
	r.Handle("/organizations", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.CreateNewOrganization))).Methods("POST")
	r.Handle("/organizations/{globalid}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:member", "organization:owner"}).Handler).Then(http.HandlerFunc(i.GetOrganization))).Methods("GET")
	r.Handle("/organizations/{globalid}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.CreateNewSubOrganization))).Methods("POST")
	r.Handle("/organizations/{globalid}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.UpdateOrganization))).Methods("PUT")
	r.Handle("/organizations/{globalid}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.DeleteOrganization))).Methods("DELETE")
	r.Handle("/organizations/{globalid}/apikeys", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.GetAPIKeyLabels))).Methods("GET")
	r.Handle("/organizations/{globalid}/apikeys", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.CreateNewAPIKey))).Methods("POST")
	r.Handle("/organizations/{globalid}/apikeys/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.GetAPIKey))).Methods("GET")
	r.Handle("/organizations/{globalid}/apikeys/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.UpdateAPIKey))).Methods("PUT")
	r.Handle("/organizations/{globalid}/apikeys/{label}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.DeleteAPIKey))).Methods("DELETE")
	r.Handle("/organizations/{globalid}/tree", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.GetOrganizationTree))).Methods("GET")
	r.Handle("/organizations/{globalid}/members", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.AddOrganizationMember))).Methods("POST")
	r.Handle("/organizations/{globalid}/members", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.UpdateOrganizationMemberShip))).Methods("PUT")
	r.Handle("/organizations/{globalid}/members/{username}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.RemoveOrganizationMember))).Methods("DELETE")
	r.Handle("/organizations/{globalid}/owners", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.AddOrganizationOwner))).Methods("POST")
	r.Handle("/organizations/{globalid}/owners/{username}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.RemoveOrganizationOwner))).Methods("DELETE")
	r.Handle("/organizations/{globalid}/contracts", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner", "organization:contracts:read"}).Handler).Then(http.HandlerFunc(i.GetContracts))).Methods("GET")
	r.Handle("/organizations/{globalid}/invitations", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.GetPendingInvitations))).Methods("GET")
	r.Handle("/organizations/{globalid}/invitations/{username}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.RemovePendingInvitation))).Methods("DELETE")
	r.Handle("/organizations/{globalid}/suborganizations", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.CreateNewSubOrganization))).Methods("POST")
	r.Handle("/organizations/{globalid}/dns/{dnsname}", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.CreateDns))).Methods("POST")
	r.Handle("/organizations/{globalid}/dns/{dnsname}", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.UpdateDns))).Methods("PUT")
	r.Handle("/organizations/{globalid}/dns/{dnsname}", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.DeleteDns))).Methods("DELETE")
	r.Handle("/organizations/{globalid}/tree", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.GetOrganizationTree))).Methods("GET")
	r.Handle("/organizations/{globalid}/registry", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.ListOrganizationRegistry))).Methods("GET")
	r.Handle("/organizations/{globalid}/registry", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.AddOrganizationRegistryEntry))).Methods("POST")
	r.Handle("/organizations/{globalid}/registry/{key}", alice.New(newOauth2oauth_2_0Middleware([]string{}).Handler).Then(http.HandlerFunc(i.GetOrganizationRegistryEntry))).Methods("GET")
	r.Handle("/organizations/{globalid}/registry/{key}", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.DeleteOrganizationRegistryEntry))).Methods("DELETE")
	r.Handle("/organizations/{globalid}/logo", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.SetOrganizationLogo))).Methods("PUT")
	r.Handle("/organizations/{globalid}/logo", http.HandlerFunc(i.GetOrganizationLogo)).Methods("GET")
	r.Handle("/organizations/{globalid}/logo", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.DeleteOrganizationLogo))).Methods("DELETE")
	r.Handle("/organizations/{globalid}/2fa/validity", http.HandlerFunc(i.Get2faValidityTime)).Methods("GET")
	r.Handle("/organizations/{globalid}/2fa/validity", alice.New(newOauth2oauth_2_0Middleware([]string{"organization:owner"}).Handler).Then(http.HandlerFunc(i.Set2faValidityTime))).Methods("PUT")
}
