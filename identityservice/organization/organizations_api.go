package organization

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"sort"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/identityservice/invitations"
	"github.com/itsyouonline/identityserver/identityservice/user"
	"github.com/itsyouonline/identityserver/oauthservice"
	"gopkg.in/mgo.v2/bson"
)

const itsyouonlineGlobalID = "itsyouonline"

// OrganizationsAPI is the implementation for /organizations root endpoint
type OrganizationsAPI struct {
}

// byGlobalID implements sort.Interface for []Organization based on
// the GlobalID field.
type byGlobalID []Organization

func (a byGlobalID) Len() int           { return len(a) }
func (a byGlobalID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byGlobalID) Less(i, j int) bool { return a[i].Globalid < a[j].Globalid }

// GetOrganizationTree is the handler for GET /organizations/{globalid}/tree
// Get organization tree.
func (api OrganizationsAPI) GetOrganizationTree(w http.ResponseWriter, r *http.Request) {
	var requestedOrganization = mux.Vars(r)["globalid"]
	//TODO: validate input
	parentGlobalID := ""
	var parentGlobalIDs = make([]string, 0, 1)
	for _, localParentID := range strings.Split(requestedOrganization, ".") {
		if parentGlobalID == "" {
			parentGlobalID = localParentID
		} else {
			parentGlobalID = parentGlobalID + "." + localParentID
		}

		parentGlobalIDs = append(parentGlobalIDs, parentGlobalID)
	}

	orgMgr := NewManager(r)

	parentOrganizations, err := orgMgr.GetOrganizations(parentGlobalIDs)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	suborganizations, err := orgMgr.GetSubOrganizations(requestedOrganization)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	allOrganizations := append(parentOrganizations, suborganizations...)

	sort.Sort(byGlobalID(allOrganizations))

	//Build a treestructure
	var orgTree *OrganizationTreeItem
	orgTreeIndex := make(map[string]*OrganizationTreeItem)
	for _, org := range allOrganizations {
		newTreeItem := &OrganizationTreeItem{GlobalID: org.Globalid, Children: make([]*OrganizationTreeItem, 0, 0)}
		orgTreeIndex[org.Globalid] = newTreeItem
		if orgTree == nil {
			orgTree = newTreeItem
		} else {
			path := strings.Split(org.Globalid, ".")
			localName := path[len(path)-1]
			parentTreeItem := orgTreeIndex[strings.TrimSuffix(org.Globalid, "."+localName)]
			parentTreeItem.Children = append(parentTreeItem.Children, newTreeItem)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orgTree)
}

// CreateNewOrganization is the handler for POST /organizations
// Create a new organization. 1 user should be in the owners list. Validation is performed
// to check if the securityScheme allows management on this user.
func (api OrganizationsAPI) CreateNewOrganization(w http.ResponseWriter, r *http.Request) {
	var org Organization

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		log.Debug("Error decoding the organization:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if strings.Contains(org.Globalid, ".") {
		log.Debug("globalid contains a '.'")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	api.actualOrganizationCreation(org, w, r)

}

// CreateNewSubOrganization is the handler for POST /organizations/{globalid}/suborganizations
// Create a new suborganization.
func (api OrganizationsAPI) CreateNewSubOrganization(w http.ResponseWriter, r *http.Request) {
	parent := mux.Vars(r)["globalid"]
	var org Organization

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		log.Debug("Error decoding the organization:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(org.Globalid, parent+".") {
		log.Debug("GlobalID does not start with the parent globalID")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	localid := strings.TrimPrefix(org.Globalid, parent+".")
	if strings.Contains(localid, ".") {
		log.Debug("localid contains a '.'")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	api.actualOrganizationCreation(org, w, r)

}

func (api OrganizationsAPI) actualOrganizationCreation(org Organization, w http.ResponseWriter, r *http.Request) {

	if strings.TrimSpace(org.Globalid) == itsyouonlineGlobalID {
		log.Debug("Duplicate organization")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if !org.IsValid() {
		log.Debug("Invalid organization")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)

	err := orgMgr.Create(&org)

	if err != nil && err != db.ErrDuplicate {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == db.ErrDuplicate {
		log.Debug("Duplicate organization")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(&org)
}

// Get organization info
// It is handler for GET /organizations/{globalid}
func (api OrganizationsAPI) globalidGet(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	orgMgr := NewManager(r)

	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(org)
}

// Update organization info
// It is handler for PUT /organizations/{globalid}
func (api OrganizationsAPI) globalidPut(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	var org Organization

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)

	oldOrg, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if org.Globalid != globalid {
		http.Error(w, "Changing globalid or id is Forbidden!", http.StatusForbidden)
		return
	}

	// Update only certain fields
	oldOrg.PublicKeys = org.PublicKeys
	oldOrg.DNS = org.DNS

	if err := orgMgr.Save(oldOrg); err != nil {
		log.Error("Error while saving organization: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oldOrg)
}

// Assign a member to organization
// It is handler for POST /organizations/{globalid}/members
func (api OrganizationsAPI) globalidmembersPost(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	var m member

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)

	org, err := orgMgr.GetByName(globalid)
	if err != nil { //TODO: make a distinction with an internal server error
		log.Debug("Error while getting the organization: ", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Check if user exists
	userMgr := user.NewManager(r)

	if ok, err := userMgr.Exists(m.Username); err != nil || !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	for _, membername := range org.Members {
		if membername == m.Username {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	for _, membername := range org.Owners {
		if membername == m.Username {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	// Create JoinRequest
	invitationMgr := invitations.NewInvitationManager(r)

	orgReq := &invitations.JoinOrganizationInvitation{
		Role:         invitations.RoleMember,
		Organization: globalid,
		User:         m.Username,
		Status:       invitations.RequestPending,
		Created:      bson.Now(),
	}

	if err := invitationMgr.Save(orgReq); err != nil {
		log.Error("Error inviting member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(orgReq)
}

// Remove a member from organization
// It is handler for DELETE /organizations/{globalid}/members/{username}
func (api OrganizationsAPI) globalidmembersusernameDelete(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	username := mux.Vars(r)["username"]

	orgMgr := NewManager(r)

	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := orgMgr.RemoveMember(org, username); err != nil {
		log.Error("Error adding member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// It is handler for POST /organizations/{globalid}/members
func (api OrganizationsAPI) globalidownersPost(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	var m member

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)

	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Check if user exists
	userMgr := user.NewManager(r)

	if ok, err := userMgr.Exists(m.Username); err != nil || !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	for _, membername := range org.Owners {
		if membername == m.Username {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}
	// Create JoinRequest
	invitationMgr := invitations.NewInvitationManager(r)

	orgReq := &invitations.JoinOrganizationInvitation{
		Role:         invitations.RoleOwner,
		Organization: globalid,
		User:         m.Username,
		Status:       invitations.RequestPending,
		Created:      bson.Now(),
	}

	if err := invitationMgr.Save(orgReq); err != nil {
		log.Error("Error inviting owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(orgReq)
}

// Remove a member from organization
// It is handler for DELETE /organizations/{globalid}/members/{username}
func (api OrganizationsAPI) globalidownersusernameDelete(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	username := mux.Vars(r)["username"]

	orgMgr := NewManager(r)

	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := orgMgr.RemoveOwner(org, username); err != nil {
		log.Error("Error removing owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPendingInvitations is the handler for GET /organizations/{globalid}/invitations
// Get the list of pending invitations for users to join this organization.
func (api OrganizationsAPI) GetPendingInvitations(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	invitationMgr := invitations.NewInvitationManager(r)

	requests, err := invitationMgr.GetPendingByOrganization(globalid)

	if err != nil {
		log.Error("Error in GetPendingByOrganization: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	pendingInvites := make([]Invitation, len(requests), len(requests))
	for index, request := range requests {
		pendingInvites[index] = Invitation{
			Role: request.Role,
			User: request.User,
			Created: request.Created,
		}
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(pendingInvites)
}

// RemovePendingInvitation is the handler for DELETE /organizations/{globalid}/invitations/{username}
// Cancel a pending invitation.
func (api OrganizationsAPI) RemovePendingInvitation(w http.ResponseWriter, r *http.Request) {
	log.Error("RemovePendingInvitation is not implemented")
}

// GetContracts is the handler for GET /organizations/{globalid}/contracts
// Get the contracts where the organization is 1 of the parties. Order descending by
// date.
func (api OrganizationsAPI) GetContracts(w http.ResponseWriter, r *http.Request) {
	log.Error("GetContracts is not implemented")
}

// GetAPISecretLabels gets the list of labels that are defined for active api secrets. The secrets themselves
// are not included.
// It is handler for GET /organizations/{globalid}/apisecrets
func (api OrganizationsAPI) GetAPISecretLabels(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]

	mgr := oauthservice.NewManager(r)
	labels, err := mgr.GetClientSecretLabels(organization)
	if err != nil {
		log.Error("Error getting a client secret labels: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(labels)
}

func isValidAPISecretLabel(label string) (valid bool) {
	valid = true
	labelLength := len(label)
	valid = valid && labelLength > 2 && labelLength < 51
	return valid
}

func isValidDNSName(label string) (valid bool) {
	valid = true
	labelLength := len(label)
	valid = valid && labelLength > 2 && labelLength < 250
	return valid
}

// GetAPISecret is the handler for GET /organizations/{globalid}/apisecrets/{label}
func (api OrganizationsAPI) GetAPISecret(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]
	label := mux.Vars(r)["label"]

	mgr := oauthservice.NewManager(r)
	secret, err := mgr.GetClientSecret(organization, label)
	if err != nil {
		log.Error("Error getting a client secret: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if secret == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	response := struct {
		Label  string `json:"label"`
		Secret string `json:"secret"`
	}{
		Label:  label,
		Secret: secret,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

// CreateNewAPISecret creates a new API Secret
// It is handler for POST /organizations/{globalid}/apisecrets
func (api OrganizationsAPI) CreateNewAPISecret(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]

	body := struct {
		Label string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !isValidAPISecretLabel(body.Label) {
		log.Debug("Invalid label: ", body.Label)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	c := oauthservice.NewOauth2Client(organization, body.Label)

	mgr := oauthservice.NewManager(r)
	err := mgr.CreateClientSecret(c)

	if err != nil && err != db.ErrDuplicate {
		log.Error("Error creating api secret label", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == db.ErrDuplicate {
		log.Debug("Duplicate label")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	response := struct {
		Label  string `json:"label"`
		Secret string `json:"secret"`
	}{
		Label:  c.Label,
		Secret: c.Secret,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

// UpdateAPISecretLabel updates the label of the secret
// It is handler for PUT /organizations/{globalid}/apisecrets/{label}
func (api OrganizationsAPI) UpdateAPISecretLabel(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]
	oldlabel := mux.Vars(r)["label"]

	log.Debug("Updating APISecret ", oldlabel)
	body := struct{ Label string }{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !isValidAPISecretLabel(body.Label) {
		log.Debug("Invalid label: ", body.Label)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mgr := oauthservice.NewManager(r)
	err := mgr.RenameClientSecret(organization, oldlabel, body.Label)

	if err != nil && err != db.ErrDuplicate {
		log.Error("Error renaming api secret label", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == db.ErrDuplicate {
		log.Debug("Duplicate label")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteAPISecret removes an API secret
// It is handler for DELETE /organizations/{globalid}/apisecrets/{label}
func (api OrganizationsAPI) DeleteAPISecret(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]
	label := mux.Vars(r)["label"]

	mgr := oauthservice.NewManager(r)
	mgr.DeleteClientSecret(organization, label)

	w.WriteHeader(http.StatusNoContent)
}

func (api OrganizationsAPI) CreateDns(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	dnsName := mux.Vars(r)["dnsname"]

	if !isValidDNSName(dnsName) {
		log.Debug("Invalid DNS name: ", dnsName)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	orgMgr := NewManager(r)
	organization, err := orgMgr.GetByName(globalid)

	err = orgMgr.AddDNS(organization, dnsName)

	if err != nil {
		log.Error("Error creating DNS name", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := struct {
		Name string `json:"name"`
	}{
		Name:  dnsName,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (api OrganizationsAPI) UpdateDns(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	oldDns := mux.Vars(r)["dnsname"]

	body := struct {
		Name string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !isValidDNSName(body.Name) {
		log.Debug("Invalid DNS name: ", body.Name)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)
	organization, err := orgMgr.GetByName(globalid)

	err = orgMgr.UpdateDNS(organization, oldDns, body.Name)

	if err != nil {
		log.Error("Error updating DNS name", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := struct {
		Name string `json:"name"`
	}{
		Name: body.Name,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (api OrganizationsAPI) DeleteDns(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	dnsName := mux.Vars(r)["dnsname"]

	orgMgr := NewManager(r)
	organization, err := orgMgr.GetByName(globalid)

	sort.Strings(organization.DNS)
	if sort.SearchStrings(organization.DNS, dnsName) == len(organization.DNS) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err = orgMgr.RemoveDNS(organization, dnsName)

	if err != nil {
		log.Error("Error removing DNS name", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNoContent)
}

