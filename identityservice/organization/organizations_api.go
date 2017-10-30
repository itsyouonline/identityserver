package organization

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"sort"

	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/gorilla/context"
	"github.com/itsyouonline/identityserver/db"
	contractdb "github.com/itsyouonline/identityserver/db/contract"
	"github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/registry"
	"github.com/itsyouonline/identityserver/db/user"
	validationdb "github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/identityservice/contract"
	"github.com/itsyouonline/identityserver/identityservice/invitations"
	"github.com/itsyouonline/identityserver/oauthservice"
	"github.com/itsyouonline/identityserver/validation"
	"gopkg.in/mgo.v2"
)

const (
	itsyouonlineGlobalID                      = "itsyouonline"
	maximumNumberOfOrganizationsPerUser       = 1000
	maximumNumberOfInvitationsPerOrganization = 10000
	DefaultLanguage                           = "en"
)

// OrganizationsAPI is the implementation for /organizations root endpoint
type OrganizationsAPI struct {
	PhonenumberValidationService  *validation.IYOPhonenumberValidationService
	EmailAddressValidationService *validation.IYOEmailAddressValidationService
}

// byGlobalID implements sort.Interface for []Organization based on
// the GlobalID field.
type byGlobalID []organization.Organization

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

	orgMgr := organization.NewManager(r)

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
	var org organization.Organization

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		log.Debug("Error decoding the organization:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	org.Globalid = strings.Trim(org.Globalid, " ")
	if !org.IsValid() {
		log.Debug("Invalid organization")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	api.actualOrganizationCreation(org, w, r)

}

// CreateNewSubOrganization is the handler for POST /organizations/{globalid}
// Create a new suborganization.
func (api OrganizationsAPI) CreateNewSubOrganization(w http.ResponseWriter, r *http.Request) {
	parent := mux.Vars(r)["globalid"]
	var org organization.Organization

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		log.Debug("Error decoding the organization:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	org.Globalid = strings.Trim(org.Globalid, " ")
	if !org.IsValidSubOrganization() {
		log.Debug("Invalid suborganization")
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

	orgMgr := organization.NewManager(r)
	if !orgMgr.Exists(parent) {
		log.Debug("Trying to create a suborganization of an unexisting parent")
		writeErrorResponse(w, http.StatusNotFound, "Parent organization does not exist")
		return
	}

	api.actualOrganizationCreation(org, w, r)

}

func (api OrganizationsAPI) actualOrganizationCreation(org organization.Organization, w http.ResponseWriter, r *http.Request) {

	if strings.TrimSpace(org.Globalid) == itsyouonlineGlobalID {
		log.Debug("Duplicate organization")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	//Clear any possible unauthorized links to organizations/users
	org.Owners = []string{}
	org.Members = []string{}
	org.OrgMembers = []string{}
	org.OrgOwners = []string{}

	orgMgr := organization.NewManager(r)
	userMgr := user.NewManager(r)
	logoMgr := organization.NewLogoManager(r)
	username := context.Get(r, "authenticateduser").(string)
	if username != "" {
		count, err := orgMgr.CountByUser(username)
		if handleServerError(w, "counting organizations by user", err) {
			return
		}
		if count >= maximumNumberOfOrganizationsPerUser {
			log.Error("Reached organization limit for user ", username)
			writeErrorResponse(w, 422, "maximum_amount_of_organizations_reached")
			return
		}
		//Set the logged in user as owner of the new organization
		org.Owners = []string{username}
	}
	userExists, err := userMgr.Exists(org.Globalid)
	if handleServerError(w, "chekcing if user exists", err) {
		return
	}
	if userExists {
		log.Debugf("Cannot create organizatino with globalid %s because a user with this name exists", org.Globalid)
		writeErrorResponse(w, http.StatusConflict, "user_exists")
		return
	}

	err = orgMgr.Create(&org)

	if err == db.ErrDuplicate {
		log.Debug("Duplicate organization")
		writeErrorResponse(w, http.StatusConflict, "organization_exists")
		return
	}
	if handleServerError(w, "creating organization", err) {
		return
	}
	err = logoMgr.Create(&org)

	if err != nil && err != db.ErrDuplicate {
		handleServerError(w, "creating organization logo", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(&org)
}

// GetOrganization Get organization info
// It is handler for GET /organizations/{globalid}
func (api OrganizationsAPI) GetOrganization(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	orgMgr := organization.NewManager(r)

	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}

	// replace owners and members usernames
	orgView, err := org.ConvertToView(user.NewManager(r), validationdb.NewManager(r))
	if handleServerError(w, "converting organization to organization view", err) {
		return
	}
	json.NewEncoder(w).Encode(orgView)
}

func (api OrganizationsAPI) inviteUser(w http.ResponseWriter, r *http.Request, role string) {
	globalID := mux.Vars(r)["globalid"]
	invitenotification := r.FormValue("invitenotification")
	if invitenotification == "" {
		invitenotification = "default"
	}
	var s searchMember

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	searchString := s.SearchString

	orgMgr := organization.NewManager(r)
	isEmailAddress := user.ValidateEmailAddress(searchString)
	isPhoneNumber := user.ValidatePhoneNumber(searchString)
	org, err := orgMgr.GetByName(globalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			writeErrorResponse(w, http.StatusNotFound, "organization_not_found")
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}

	u, err := SearchUser(r, searchString)
	if err == mgo.ErrNotFound {
		if !isEmailAddress && !isPhoneNumber {
			writeErrorResponse(w, http.StatusNotFound, "user_not_found")
			return
		}
	} else if handleServerError(w, "searching for user", err) {
		return
	}
	username := ""
	emailAddress := ""
	code := ""
	phoneNumber := ""
	method := invitations.MethodWebsite
	if u == nil {
		randombytes := make([]byte, 9) //Multiple of 3 to make sure no padding is added
		rand.Read(randombytes)
		code = base64.URLEncoding.EncodeToString(randombytes)
		if isEmailAddress {
			emailAddress = searchString
			method = invitations.MethodEmail
		} else if isPhoneNumber {
			phoneNumber = searchString
			method = invitations.MethodPhone
		}
	} else {
		username = u.Username
		if role == invitations.RoleMember {
			for _, membername := range org.Members {
				if membername == u.Username {
					http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
					return
				}
			}
		}
		for _, memberName := range org.Owners {
			if memberName == username {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}
		}
	}
	// Create JoinRequest
	invitationMgr := invitations.NewInvitationManager(r)
	count, err := invitationMgr.CountByOrganization(globalID)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if count >= maximumNumberOfInvitationsPerOrganization {
		log.Error("Reached invitation limit for organization ", globalID)
		writeErrorResponse(w, 422, "max_amount_of_invitations_reached")
		return
	}

	orgReq := &invitations.JoinOrganizationInvitation{
		Role:           role,
		Organization:   globalID,
		User:           username,
		Status:         invitations.RequestPending,
		Created:        db.DateTime(time.Now()),
		Method:         method,
		EmailAddress:   emailAddress,
		PhoneNumber:    phoneNumber,
		Code:           code,
		IsOrganization: false,
	}

	autoAccepted, err := api.autoAcceptThreefoldInviteIfPossible(orgReq, orgMgr)
	if err != nil {
		log.Error("Failure while trying to auto accept organization invite: ", err)
		log.Warnf("OrgId: %s, role: %s")
		if db.IsNotFound(err) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err = invitationMgr.Save(orgReq); err != nil {
		log.Error("Error inviting owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if invitenotification != "none" && !autoAccepted {
		err = api.sendInvite(r, orgReq)
		if handleServerError(w, "sending organization invite", err) {
			return
		}
	}

	usrMgr := user.NewManager(r)
	valMgr := validationdb.NewManager(r)
	reqView, err := orgReq.ConvertToView(usrMgr, valMgr)
	if handleServerError(w, "converting invite to inviteview", err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reqView)
}

// FIXME: NUKE ASAP
func (api OrganizationsAPI) autoAcceptThreefoldInviteIfPossible(invite *invitations.JoinOrganizationInvitation, mgr *organization.Manager) (bool, error) {
	globalid := invite.Organization

	// we can't just check that the prefix is "threefold" as this would include unrelated
	// organizations such as threefoldwithsomethingbehindit
	if !(globalid == "threefold" || strings.HasPrefix(globalid, "threefold.")) {
		return false, nil
	}

	// Check for username so we only auto invite already registered users
	username := invite.User
	if username == "" {
		log.Debug("can't auto accept the invite because the username is not known")
		return false, nil
	}
	// add the user
	log.Debug("Try auto adding user to organization ", invite.Organization)
	org, err := mgr.GetByName(invite.Organization)
	if err != nil {
		return false, err
	}
	if invitations.RoleOwner == invite.Role {
		// Accepted Owner role
		if err := mgr.SaveOwner(org, username); err != nil {
			return false, err
		}
	} else {
		// Accepted member role
		if err := mgr.SaveMember(org, username); err != nil {
			return false, err
		}
	}
	log.Debug("Auto added user ", username, " to the ", globalid, " organization")
	// set the invite to status accpeted
	invite.Status = invitations.RequestAccepted
	return true, nil
}

// AddOrganizationMember Assign a member to organization
// It is handler for POST /organizations/{globalid}/members
func (api OrganizationsAPI) AddOrganizationMember(w http.ResponseWriter, r *http.Request) {
	api.inviteUser(w, r, invitations.RoleMember)
}

func (api OrganizationsAPI) UpdateOrganizationMemberShip(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	var membership Membership
	if err := json.NewDecoder(r.Body).Decode(&membership); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	orgMgr := organization.NewManager(r)
	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			handleServerError(w, "updating organization membership", err)
		}
		return
	}
	var oldRole string
	var username string
	valMgr := validationdb.NewManager(r)
	members, err := organization.MapUsernamesToIdentifiers(org.Members, valMgr)
	if handleServerError(w, "Converting usernames to identifiers", err) {
		return
	}
	for identifier, uname := range members {
		if identifier == membership.Username {
			oldRole = "members"
			username = uname
			break
		}
	}
	owners, err := organization.MapUsernamesToIdentifiers(org.Owners, valMgr)
	if handleServerError(w, "Converting usernames to identifiers", err) {
		return
	}
	for identifier, uname := range owners {
		if identifier == membership.Username {
			oldRole = "owners"
			username = uname
			break
		}
	}
	// if oldRole is still not set then the given user is not part of this organization
	if oldRole == "" {
		writeErrorResponse(w, http.StatusNotFound, "user_not_member_of_organization")
		return
	}
	err = orgMgr.UpdateMembership(globalid, username, oldRole, membership.Role)
	if err != nil {
		handleServerError(w, "updating organization membership", err)
		return
	}
	org, err = orgMgr.GetByName(globalid)
	if err != nil {
		handleServerError(w, "getting organization", err)
	}
	usrMgr := user.NewManager(r)
	orgView, err := org.ConvertToView(usrMgr, valMgr)
	if handleServerError(w, "converting organization to view", err) {
		return
	}
	json.NewEncoder(w).Encode(orgView)

}

func (api OrganizationsAPI) UpdateOrganizationOrgMemberShip(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	body := struct {
		Org  string
		Role string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := organization.NewManager(r)
	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			handleServerError(w, "updating organization membership", err)
		}
		return
	}

	if !orgMgr.Exists(body.Org) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var oldRole string
	for _, v := range org.OrgMembers {
		if v == body.Org {
			oldRole = "orgmembers"
		}
	}
	for _, v := range org.OrgOwners {
		if v == body.Org {
			oldRole = "orgowners"
		}
	}
	if body.Role == "members" {
		body.Role = "orgmembers"
	} else {
		body.Role = "orgowners"
	}
	err = orgMgr.UpdateOrgMembership(globalid, body.Org, oldRole, body.Role)
	if handleServerError(w, "updating organizations membership in another org", err) {
		return
	}
	org, err = orgMgr.GetByName(globalid)
	if handleServerError(w, "getting organization", err) {
		return
	}
	json.NewEncoder(w).Encode(org)

}

// RemoveOrganizationMember Remove a member from organization
// It is handler for DELETE /organizations/{globalid}/members/{username}
func (api OrganizationsAPI) RemoveOrganizationMember(w http.ResponseWriter, r *http.Request) {
	removeOrganizationMember(w, r, "member")
}

func removeOrganizationMember(w http.ResponseWriter, r *http.Request, role string) {
	globalID := mux.Vars(r)["globalid"]
	userIdentifier := mux.Vars(r)["username"]

	orgMgr := organization.NewManager(r)
	userMgr := user.NewManager(r)
	valMgr := validationdb.NewManager(r)
	username, err := organization.ConvertIdentifierToUsername(userIdentifier, valMgr)
	if handleServerError(w, "Converting username to identifier", err) {
		return
	}
	org, err := orgMgr.GetByName(globalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}
	if role == "member" {
		if handleServerError(w, "Removing organization member", orgMgr.RemoveMember(org, username)) {
			return
		}
	} else if role == "owner" {
		if handleServerError(w, "Removing organization owner", orgMgr.RemoveOwner(org, username)) {
			return
		}
	} else {
		log.Errorf("Invalid role given to removeOrganizationMember: %s", role)
		writeErrorResponse(w, http.StatusInternalServerError, "invalid_role")
	}

	err = userMgr.DeleteAuthorization(username, globalID)
	if handleServerError(w, "removing authorization", err) {
		return
	}

	invitationMgr := invitations.NewInvitationManager(r)
	err = invitationMgr.Remove(globalID, username, username)
	if db.IsNotFound(err) {
		// most of the time the users will have no invitation if they are already part
		// of the organization so just silently ignore this
		err = nil
	}
	if handleServerError(w, "removing invitation", err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddOrganizationOwner It is handler for POST /organizations/{globalid}/owners
func (api OrganizationsAPI) AddOrganizationOwner(w http.ResponseWriter, r *http.Request) {
	api.inviteUser(w, r, invitations.RoleOwner)
}

func (api OrganizationsAPI) sendInvite(r *http.Request, organizationRequest *invitations.JoinOrganizationInvitation) error {
	switch organizationRequest.Method {
	case invitations.MethodWebsite:
		return nil
	case invitations.MethodEmail:
		return api.EmailAddressValidationService.SendOrganizationInviteEmail(r, organizationRequest)
	case invitations.MethodPhone:
		return api.PhonenumberValidationService.SendOrganizationInviteSms(r, organizationRequest)
	}
	return nil
}

// RemoveOrganizationOwner Remove a member from organization
// It is handler for DELETE /organizations/{globalid}/owners/{username}
func (api OrganizationsAPI) RemoveOrganizationOwner(w http.ResponseWriter, r *http.Request) {
	removeOrganizationMember(w, r, "owner")
}

// GetInvitations is the handler for GET /organizations/{globalid}/invitations
// Get the list of pending invitations for users to join this organization.
func (api OrganizationsAPI) GetInvitations(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	status := invitations.ParseInvitationType(r.FormValue("status"))

	invitationMgr := invitations.NewInvitationManager(r)
	invites, err := invitationMgr.FilterByOrganization(globalid, status)

	if handleServerError(w, "filtering invitations by organization", err) {
		return
	}

	usrMgr := user.NewManager(r)
	valMgr := validationdb.NewManager(r)
	views := make([]*invitations.JoinOrganizationInvitationView, len(invites))
	for i, invite := range invites {
		// todo: query in loop -> will be slow for lots of invites
		view, err := invite.ConvertToView(usrMgr, valMgr)
		if handleServerError(w, "converting invite to view", err) {
			return
		}
		views[i] = view
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(views)
}

// RemovePendingInvitation is the handler for DELETE /organizations/{globalid}/invitations/{username}
// Cancel a pending invitation.
func (api OrganizationsAPI) RemovePendingInvitation(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]
	searchString := mux.Vars(r)["searchstring"]
	invitationMgr := invitations.NewInvitationManager(r)

	usr, err := SearchUser(r, searchString)
	if err != nil {
		if db.IsNotFound(err) {
			log.Debug("Could not find existing user for whom to remove the invitation")
			usr = &user.User{}
			// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			// return
		} else {
			log.Error("Error while searching user: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	err = invitationMgr.Remove(globalID, usr.Username, searchString)
	if err == mgo.ErrNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error("Error while remove invite: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetContracts is the handler for GET /organizations/{globalid}/contracts
// Get the contracts where the organization is 1 of the parties. Order descending by
// date.
func (api OrganizationsAPI) GetContracts(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalId"]
	includedparty := contractdb.Party{Type: "org", Name: globalID}
	contract.FindContracts(w, r, includedparty)
}

// RegisterNewContract is handler for POST /organizations/{globalId}/contracts
func (api OrganizationsAPI) RegisterNewContract(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["glabalId"]
	includedparty := contractdb.Party{Type: "org", Name: globalID}
	contract.CreateContract(w, r, includedparty)
}

// GetAPIKeyLabels is the handler for GET /organizations/{globalid}/apikeys
// Get the list of active api keys. The secrets themselves are not included.
func (api OrganizationsAPI) GetAPIKeyLabels(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]

	mgr := oauthservice.NewManager(r)
	labels, err := mgr.GetClientLabels(organization)
	if err != nil {
		log.Error("Error getting a client secret labels: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(labels)
}

// GetAPIKey is the handler for GET /organizations/{globalid}/apikeys/{label}
func (api OrganizationsAPI) GetAPIKey(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]
	label := mux.Vars(r)["label"]

	mgr := oauthservice.NewManager(r)
	client, err := mgr.GetClient(organization, label)
	if err != nil {
		log.Error("Error getting a client: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if client == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	apiKey := FromOAuthClient(client)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(apiKey)
}

// CreateNewAPIKey is the handler for POST /organizations/{globalid}/apikeys
// Create a new API Key, a secret itself should not be provided, it will be generated
// serverside.
func (api OrganizationsAPI) CreateNewAPIKey(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]

	apiKey := APIKey{}

	if err := json.NewDecoder(r.Body).Decode(&apiKey); err != nil {
		log.Debug("Error decoding apikey: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	apiKey.Secret = "dummysecret" // else the validator complains
	if !apiKey.Validate() {
		log.Debug("Invalid api key: ", apiKey)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	log.Debug("Creating apikey:", apiKey)
	c := oauthservice.NewOauth2Client(globalID, apiKey.Label, apiKey.CallbackURL, apiKey.ClientCredentialsGrantType)

	mgr := oauthservice.NewManager(r)
	err := mgr.CreateClient(c)
	if db.IsDup(err) {
		log.Debug("Duplicate label")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	if err != nil {
		log.Error("Error creating api secret label", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	apiKey.Secret = c.Secret

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiKey)

}

// UpdateAPIKey is the handler for PUT /organizations/{globalid}/apikeys/{label}
// Updates the label or other properties of a key.
func (api OrganizationsAPI) UpdateAPIKey(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]
	oldLabel := mux.Vars(r)["label"]

	apiKey := APIKey{}

	if err := json.NewDecoder(r.Body).Decode(&apiKey); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// set a fake secret or override the existing secret for the validator. The secret
	// is ignored anyway when updating the apikey
	apiKey.Secret = "dummysecret"
	if !apiKey.Validate() {
		log.Debug("Invalid api key: ", apiKey)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mgr := oauthservice.NewManager(r)
	c, err := mgr.GetClient(globalID, oldLabel)
	if handleServerError(w, "getting old api key", err) {
		return
	}
	if c == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	err = mgr.UpdateClient(globalID, oldLabel, apiKey.Label, apiKey.CallbackURL, apiKey.ClientCredentialsGrantType)

	if err != nil && db.IsDup(err) {
		log.Debug("Duplicate label")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if err != nil {
		log.Error("Error renaming api secret label", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteAPIKey is the handler for DELETE /organizations/{globalid}/apikeys/{label}
// Removes an API key
func (api OrganizationsAPI) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]
	label := mux.Vars(r)["label"]

	mgr := oauthservice.NewManager(r)
	err := mgr.DeleteClient(organization, label)

	if err != nil {
		log.Error("Error deleting organization:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CreateOrganizationDns is the handler for POST /organizations/{globalid}/dns
// Adds a dns address to an organization
func (api OrganizationsAPI) CreateOrganizationDns(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]

	dns := DnsAddress{}

	if err := json.NewDecoder(r.Body).Decode(&dns); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !dns.Validate() {
		log.Debug("Invalid DNS name: ", dns.Name)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	orgMgr := organization.NewManager(r)
	organisation, err := orgMgr.GetByName(globalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}
	for _, d := range organisation.DNS {
		if dns.Name == d {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}
	err = orgMgr.AddDNS(organisation, dns.Name)

	if handleServerError(w, "adding DNS name", err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dns)
}

// UpdateOrganizationDns is the handler for PUT /organizations/{globalid}/dns/{dnsname}
// Updates an existing DNS name associated with an organization
func (api OrganizationsAPI) UpdateOrganizationDns(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]
	oldDNS := mux.Vars(r)["dnsname"]

	var dns DnsAddress

	if err := json.NewDecoder(r.Body).Decode(&dns); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !dns.Validate() {
		log.Debug("Invalid DNS name: ", dns.Name)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := organization.NewManager(r)
	organisation, err := orgMgr.GetByName(globalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}
	exists := false
	for _, d := range organisation.DNS {
		if d == dns.Name {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
		if d == oldDNS {
			exists = true
		}
	}
	if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	err = orgMgr.UpdateDNS(organisation, oldDNS, dns.Name)

	if err != nil {
		log.Error("Error updating DNS name", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dns)
}

// DeleteOrganizationDns is the handler for DELETE /organizations/{globalid}/dns/{dnsname}
// Removes a DNS name associated with an organization
func (api OrganizationsAPI) DeleteOrganizationDns(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	dnsName := mux.Vars(r)["dnsname"]

	orgMgr := organization.NewManager(r)
	organisation, err := orgMgr.GetByName(globalid)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}
	sort.Strings(organisation.DNS)
	if sort.SearchStrings(organisation.DNS, dnsName) == len(organisation.DNS) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err = orgMgr.RemoveDNS(organisation, dnsName)

	if err != nil {
		log.Error("Error removing DNS name", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNoContent)
}

// DeleteOrganization is the handler for DELETE /organizations/{globalid}
// Deletes an organization and all data linked to it (join-organization-invitations, oauth_access_tokens, oauth_clients, authorizations)
func (api OrganizationsAPI) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	orgMgr := organization.NewManager(r)
	suborganizations, err := orgMgr.GetSubOrganizations(globalid)
	if handleServerError(w, "fetching suborganizations", err) {
		return
	}
	for _, org := range suborganizations {
		api.actualOrganizationDeletion(w, r, org.Globalid)
	}
	api.actualOrganizationDeletion(w, r, globalid)
}

// Delete organization with globalid.
func (api OrganizationsAPI) actualOrganizationDeletion(w http.ResponseWriter, r *http.Request, globalid string) {
	orgMgr := organization.NewManager(r)
	logoMgr := organization.NewLogoManager(r)
	if !orgMgr.Exists(globalid) {
		writeErrorResponse(w, http.StatusNotFound, "organization_not_found")
		return
	}
	err := orgMgr.Remove(globalid)
	if handleServerError(w, "removing organization", err) {
		return
	}
	// Remove the organizations as a member/ an owner of other organizations
	organizations, err := orgMgr.AllByOrg(globalid)
	if handleServerError(w, "fetching organizations where this org is an owner/a member", err) {
		return
	}
	for _, org := range organizations {
		err = orgMgr.RemoveOrganization(org.Globalid, globalid)
		if handleServerError(w, "removing organizations as a member / an owner of another organization", err) {
			return
		}
	}
	if logoMgr.Exists(globalid) {
		err = logoMgr.Remove(globalid)
		if handleServerError(w, "removing organization logo", err) {
			return
		}
	}
	orgReqMgr := invitations.NewInvitationManager(r)
	err = orgReqMgr.RemoveAll(globalid)
	if handleServerError(w, "removing organization invitations", err) {
		return
	}

	oauthMgr := oauthservice.NewManager(r)
	err = oauthMgr.RemoveTokensByGlobalID(globalid)
	if handleServerError(w, "removing organization oauth accesstokens", err) {
		return
	}
	err = oauthMgr.DeleteAllForOrganization(globalid)
	if handleServerError(w, "removing client secrets", err) {
		return
	}
	err = oauthMgr.RemoveClientsByID(globalid)
	if handleServerError(w, "removing organization oauth clients", err) {
		return
	}
	userMgr := user.NewManager(r)
	err = userMgr.DeleteAllAuthorizations(globalid)
	if handleServerError(w, "removing all authorizations", err) {
		return
	}
	err = oauthMgr.RemoveClientsByID(globalid)
	if handleServerError(w, "removing organization oauth clients", err) {
		return
	}
	l2faMgr := organization.NewLast2FAManager(r)
	err = l2faMgr.RemoveByOrganization(globalid)
	if handleServerError(w, "removing organization 2FA history", err) {
		return
	}
	descriptionMgr := organization.NewDescriptionManager(r)
	err = descriptionMgr.Remove(globalid)
	if err != nil {
		if err != mgo.ErrNotFound {
			handleServerError(w, "removing organization description", err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListOrganizationRegistry is the handler for GET /organizations/{globalid}/registry
// Lists the Registry entries
func (api OrganizationsAPI) ListOrganizationRegistry(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	mgr := registry.NewManager(r)
	registryEntries, err := mgr.ListRegistryEntries("", globalid)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registryEntries)
}

// AddOrganizationRegistryEntry is the handler for POST /organizations/{globalid}/registry
// Adds a RegistryEntry to the organization's registry, if the key is already used, it is overwritten.
func (api OrganizationsAPI) AddOrganizationRegistryEntry(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	registryEntry := registry.RegistryEntry{}

	if err := json.NewDecoder(r.Body).Decode(&registryEntry); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := registryEntry.Validate(); err != nil {
		log.Debug("Invalid registry entry: ", registryEntry)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mgr := registry.NewManager(r)
	err := mgr.UpsertRegistryEntry("", globalid, registryEntry)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registryEntry)
}

// GetOrganizationRegistryEntry is the handler for GET /organizations/{username}/globalid/{key}
// Get a RegistryEntry from the organization's registry.
func (api OrganizationsAPI) GetOrganizationRegistryEntry(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	key := mux.Vars(r)["key"]

	mgr := registry.NewManager(r)
	registryEntry, err := mgr.GetRegistryEntry("", globalid, key)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if registryEntry == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registryEntry)
}

// DeleteOrganizationRegistryEntry is the handler for DELETE /organizations/{username}/globalid/{key}
// Removes a RegistryEntry from the organization's registry
func (api OrganizationsAPI) DeleteOrganizationRegistryEntry(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	key := mux.Vars(r)["key"]

	mgr := registry.NewManager(r)
	err := mgr.DeleteRegistryEntry("", globalid, key)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SetOrganizationLogo is the handler for PUT /organizations/globalid/logo
// Set the organization Logo for the organization
func (api OrganizationsAPI) SetOrganizationLogo(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	body := struct {
		Logo string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("Error while saving logo: ", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	logoMgr := organization.NewLogoManager(r)

	// server side file size validation check. Normally uploaded files should never get this large due to size constraints, but check anyway
	if len(body.Logo) > 1024*1024*5 {
		log.Error("Error while saving file: file too large")
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}
	_, err := logoMgr.SaveLogo(globalid, body.Logo)
	if err != nil {
		log.Error("Error while saving logo: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetOrganizationLogo is the handler for GET /organizations/globalid/logo
// Get the Logo from an organization
func (api OrganizationsAPI) GetOrganizationLogo(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	logoMgr := organization.NewLogoManager(r)

	logo, err := logoMgr.GetLogo(globalid)

	if err != nil && err != mgo.ErrNotFound {
		log.Error("Error getting logo", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := struct {
		Logo string `json:"logo"`
	}{
		Logo: logo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteOrganizationLogo is the handler for DELETE /organizations/globalid/logo
// Removes the Logo from an organization
func (api OrganizationsAPI) DeleteOrganizationLogo(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	logoMgr := organization.NewLogoManager(r)

	err := logoMgr.RemoveLogo(globalid)

	if err != nil {
		log.Error("Error removing logo", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNoContent)
}

// Get2faValidityTime is the handler for GET /organizations/globalid/2fa/validity
// Get the 2fa validity time for the organization, in seconds
func (api OrganizationsAPI) Get2faValidityTime(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	mgr := organization.NewManager(r)

	validity, err := mgr.GetValidity(globalid)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("Error while getting validity duration: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err == mgo.ErrNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	response := struct {
		SecondsValidity int `json:"secondsvalidity"`
	}{
		SecondsValidity: validity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

// Set2faValidityTime is the handler for PUT /organizations/globalid/2fa/validity
// Sets the 2fa validity time for the organization, in days
func (api OrganizationsAPI) Set2faValidityTime(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	body := struct {
		SecondsValidity int `json:"secondsvalidity"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("Error while setting 2FA validity time: ", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mgr := organization.NewManager(r)
	seconds := body.SecondsValidity

	if seconds < 0 {
		seconds = 0
	} else if seconds > 2678400 {
		seconds = 2678400
	}

	err := mgr.SetValidity(globalid, seconds)
	if err != nil {
		log.Error("Error while setting 2FA validity time: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SetOrgMember is the handler for POST /organizations/{globalid}/orgmember
// Sets an organization as a member of this one.
func (api OrganizationsAPI) SetOrgMember(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	body := struct {
		OrgMember string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Debug("Error while adding another organization as member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mgr := organization.NewManager(r)

	// load organization for globalid
	organization, err := mgr.GetByName(globalid)

	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}

	// check if OrgMember exists
	if !mgr.Exists(body.OrgMember) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// now that we know both organizations exists, check if the authenticated user is an owner of the OrgMember
	// the user is known to be an owner of the first organization since we've required the organization:owner scope
	authenticateduser := context.Get(r, "authenticateduser").(string)
	isOwner, err := mgr.IsOwner(body.OrgMember, authenticateduser)
	if err != nil {
		log.Error("Error while adding another organization as member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !isOwner {
		// invite the organization if we can't add it directly
		api.inviteOrganization(w, r, invitations.RoleOrgMember, body.OrgMember)
		return
	}

	// check if thie organization we want to add already exists as a member or an owner
	exists, err := mgr.OrganizationIsPartOf(globalid, body.OrgMember)
	if err != nil {
		log.Error("Error while checking if this organization is part of another: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	err = mgr.SaveOrgMember(organization, body.OrgMember)
	if err != nil {
		log.Error("Error while adding another organization as member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteOrgMember is the handler for Delete /organizations/globalid/orgmember/globalid2
// Removes an organization as a member of this one.
func (api OrganizationsAPI) DeleteOrgMember(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	orgMember := mux.Vars(r)["globalid2"]

	mgr := organization.NewManager(r)

	if !mgr.Exists(globalid) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// check if OrgMember is a member of the organization
	isMember, err := mgr.OrganizationIsMember(globalid, orgMember)
	if err != nil {
		log.Error("Error while removing another organization as member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !isMember {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err = mgr.RemoveOrganization(globalid, orgMember)
	if err != nil {
		log.Error("Error while removing another organization as member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SetOrgOwner is the handler for POST /organizations/globalid/orgowner
// Sets an organization as an owner of this one.
func (api OrganizationsAPI) SetOrgOwner(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	body := struct {
		OrgOwner string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Debug("Error while adding another organization as owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mgr := organization.NewManager(r)

	// load organization for globalid
	organization, err := mgr.GetByName(globalid)

	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}

	// check if OrgOwner exists
	if !mgr.Exists(body.OrgOwner) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// now that we know both organizations exists, check if the authenticated user is an owner of the OrgOwner
	// the user is known to be an owner of the first organization since we've required the organization:owner scope
	authenticateduser := context.Get(r, "authenticateduser").(string)
	isOwner, err := mgr.IsOwner(body.OrgOwner, authenticateduser)
	if err != nil {
		log.Error("Error while adding another organization as owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !isOwner {
		api.inviteOrganization(w, r, invitations.RoleOrgOwner, body.OrgOwner)
		return
	}

	// check if the organization we want to add already exists as a member or an owner
	exists, err := mgr.OrganizationIsPartOf(globalid, body.OrgOwner)
	if err != nil {
		log.Error("Error while checking if this organization is part of another: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	err = mgr.SaveOrgOwner(organization, body.OrgOwner)
	if err != nil {
		log.Error("Error while adding another organization as owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteOrgOwner is the handler for Delete /organizations/globalid/orgowner/globalid2
// Removes an organization as an owner of this one.
func (api OrganizationsAPI) DeleteOrgOwner(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	orgOwner := mux.Vars(r)["globalid2"]

	mgr := organization.NewManager(r)

	if !mgr.Exists(globalid) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// check if OrgOwner is an owner of the organization
	isOwner, err := mgr.OrganizationIsOwner(globalid, orgOwner)
	if err != nil {
		log.Error("Error while removing another organization as owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !isOwner {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err = mgr.RemoveOrganization(globalid, orgOwner)
	if err != nil {
		log.Error("Error while removing another organization as owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddRequiredScope is the handler for POST /organizations/{globalid}/requiredscope
// Adds a required scope
func (api OrganizationsAPI) AddRequiredScope(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]
	var requiredScope organization.RequiredScope
	if err := json.NewDecoder(r.Body).Decode(&requiredScope); err != nil {
		log.Debug("Error while adding a required scope: ", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !requiredScope.IsValid() {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	mgr := organization.NewManager(r)
	organisation, err := mgr.GetByName(globalID)
	if err == mgo.ErrNotFound {
		writeErrorResponse(w, http.StatusNotFound, "organization_not_found")
		return
	}
	for _, scope := range organisation.RequiredScopes {
		if scope.Scope == requiredScope.Scope {
			writeErrorResponse(w, http.StatusConflict, "required_scope_already_exists")
			return
		}
	}
	err = mgr.AddRequiredScope(globalID, requiredScope)
	if err != nil {
		handleServerError(w, "adding a required scope", err)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// UpdateRequiredScope is the handler for PUT /organizations/{globalid}/requiredscope/{requiredscope}
// Updates a required scope
func (api OrganizationsAPI) UpdateRequiredScope(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]
	oldRequiredScope := mux.Vars(r)["requiredscope"]
	var requiredScope organization.RequiredScope
	if err := json.NewDecoder(r.Body).Decode(&requiredScope); err != nil {
		log.Debug("Error while updating a required scope: ", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !requiredScope.IsValid() {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	mgr := organization.NewManager(r)
	exists := mgr.Exists(globalID)
	if !exists {
		writeErrorResponse(w, http.StatusNotFound, "organization_not_found")
		return
	}
	err := mgr.UpdateRequiredScope(globalID, oldRequiredScope, requiredScope)
	if err != nil {
		if err == mgo.ErrNotFound {
			writeErrorResponse(w, http.StatusNotFound, "required_scope_not_found")
		} else {
			handleServerError(w, "updating a required scope", err)
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// DeleteRequiredScope is the handler for DELETE /organizations/{globalid}/requiredscope/{requiredscope}
// Updates a required scope
func (api OrganizationsAPI) DeleteRequiredScope(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]
	requiredScope := mux.Vars(r)["requiredscope"]
	mgr := organization.NewManager(r)
	if !mgr.Exists(globalID) {
		writeErrorResponse(w, http.StatusNotFound, "organization_not_found")
		return
	}
	err := mgr.DeleteRequiredScope(globalID, requiredScope)
	if err != nil {
		if err == mgo.ErrNotFound {
			writeErrorResponse(w, http.StatusNotFound, "required_scope_not_found")
		} else {
			handleServerError(w, "removing a required scope", err)
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// GetOrganizationUsers is the handler for GET /organizations/{globalid}/users
// Get the list of all users in this organization
func (api OrganizationsAPI) GetOrganizationUsers(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]
	orgMgr := organization.NewManager(r)
	if !orgMgr.Exists(globalID) {
		writeErrorResponse(w, http.StatusNotFound, "organization_not_found")
		return
	}
	authenticatedUser := context.Get(r, "authenticateduser").(string)
	response := organization.GetOrganizationUsersResponseBody{}
	isOwner, err := orgMgr.IsOwner(globalID, authenticatedUser)
	if handleServerError(w, "checking if user is owner of an organization", err) {
		return
	}
	org, err := orgMgr.GetByName(globalID)
	if handleServerError(w, "getting organization by name", err) {
		return
	}
	roleMap := make(map[string]string)
	valMgr := validationdb.NewManager(r)
	allUsernames := append(org.Members, org.Owners...)
	userIdentifierMap, err := organization.MapUsernamesToIdentifiers(allUsernames, valMgr)
	if handleServerError(w, "mapping usernames to identifiers", err) {
		return
	}
	for _, username := range org.Members {
		roleMap[username] = "members"
	}
	for _, username := range org.Owners {
		roleMap[username] = "owners"
	}
	authorizationsMap := make(map[string]user.Authorization)
	// Only owners can see if there are missing permissions
	userMgr := user.NewManager(r)
	if isOwner {
		authorizations, err := userMgr.GetOrganizationAuthorizations(globalID)
		if handleServerError(w, "getting organizaton authorizations", err) {
			return
		}
		for _, authorization := range authorizations {
			authorizationsMap[authorization.Username] = authorization
		}
	}
	users := []organization.OrganizationUser{}
	for username, role := range roleMap {
		orgUser := organization.OrganizationUser{
			Username:      getUserIdentifier(username, userIdentifierMap),
			Role:          role,
			MissingScopes: []string{},
		}
		if isOwner {
			for _, requiredScope := range org.RequiredScopes {
				hasScope := false
				if authorization, hasKey := authorizationsMap[username]; hasKey {
					hasScope = requiredScope.IsAuthorized(authorization)
				} else {
					hasScope = false
				}
				if !hasScope {
					orgUser.MissingScopes = append(orgUser.MissingScopes, requiredScope.Scope)
				}
			}
		}
		users = append(users, orgUser)
	}
	response.HasEditPermissions = isOwner
	response.Users = users
	json.NewEncoder(w).Encode(response)
}

type sortEmailFirst []string

func (s sortEmailFirst) Len() int {
	return len(s)
}
func (s sortEmailFirst) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortEmailFirst) Less(i, j int) bool {
	if strings.Contains(s[i], "@") {
		if !strings.Contains(s[j], "@") {
			return true
		} else {
			return sort.StringsAreSorted([]string{s[i], s[j]})
		}
	}
	return false
}

func getUserIdentifier(username string, userIdentifierMap map[string]string) string {
	identifiers := []string{}
	for identifier, uname := range userIdentifierMap {
		if username == uname {
			identifiers = append(identifiers, identifier)
		}
	}
	if len(identifiers) > 0 {
		sort.Sort(sortEmailFirst(identifiers))
		return identifiers[0]
	}
	// fallback to username
	return username
}

// GetDescription is the handler for GET /organizations/{globalid}/description/{langkey}
// Get the description for this organization for this langKey
func (api OrganizationsAPI) GetDescription(w http.ResponseWriter, r *http.Request) {
	globalId := mux.Vars(r)["globalid"]
	rawLangKey := mux.Vars(r)["langkey"]
	langKey := parseLangKey(rawLangKey)
	if langKey == "" {
		writeErrorResponse(w, http.StatusBadRequest, "invalid_language_key")
		return
	}
	mgr := organization.NewDescriptionManager(r)
	var description organization.LocalizedInfoText
	orgDescriptions, err := mgr.GetDescription(globalId)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Error("ERROR while loading localized description: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		description.LangKey = langKey
		description.Text = ""
	}
	for _, storedDescription := range orgDescriptions.InfoTexts {
		if storedDescription.LangKey == langKey {
			description = storedDescription
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(description)
}

// GetDescriptionWithFallback is the handler for GET /organizations/{globalid}/description/{langkey}/withfallback
// Get the description for this organization for this langKey. If it doesn't exist, get the desription for the default langKey
func (api OrganizationsAPI) GetDescriptionWithFallback(w http.ResponseWriter, r *http.Request) {
	globalId := mux.Vars(r)["globalid"]
	rawLangKey := mux.Vars(r)["langkey"]
	langKey := parseLangKey(rawLangKey)
	if langKey == "" {
		writeErrorResponse(w, http.StatusBadRequest, "invalid_language_key")
		return
	}
	mgr := organization.NewDescriptionManager(r)
	var description organization.LocalizedInfoText
	orgDescriptions, err := mgr.GetDescription(globalId)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Error("ERROR while loading localized description: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		description.LangKey = langKey
		description.Text = ""
	}
	for _, storedDescription := range orgDescriptions.InfoTexts {
		if storedDescription.LangKey == langKey {
			description = storedDescription
		}
	}

	// If no translation is found for the langKey, try the default langKey
	if description.Text == "" && langKey != DefaultLanguage {
		langKey = DefaultLanguage
	}
	for _, storedDescription := range orgDescriptions.InfoTexts {
		if storedDescription.LangKey == langKey {
			description = storedDescription
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(description)
}

// DeleteDescription is the handler for GET /organizations/{globalid}/description/{langkey}
// Delete the description for this organization for this langKey
func (api OrganizationsAPI) DeleteDescription(w http.ResponseWriter, r *http.Request) {
	globalId := mux.Vars(r)["globalid"]
	rawLangKey := mux.Vars(r)["langkey"]
	langKey := parseLangKey(rawLangKey)
	if langKey == "" {
		writeErrorResponse(w, http.StatusBadRequest, "invalid_language_key")
		return
	}
	mgr := organization.NewDescriptionManager(r)
	err := mgr.DeleteDescription(globalId, langKey)
	if err != nil {
		log.Error("ERROR while deleting localized description: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// SetDescription is the handler for POST /organizations/{globalid}/description
// Set the description for this organization for this langKey
func (api OrganizationsAPI) SetDescription(w http.ResponseWriter, r *http.Request) {
	globalId := mux.Vars(r)["globalid"]
	var localInfo organization.LocalizedInfoText

	if err := json.NewDecoder(r.Body).Decode(&localInfo); err != nil {
		log.Debug("Error decoding the localized description:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	langKey := parseLangKey(localInfo.LangKey)
	if langKey == "" {
		log.Debug("Error decoding the localized description: invalid language key")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	localInfo.LangKey = langKey
	descriptionMgr := organization.NewDescriptionManager(r)
	err := descriptionMgr.SaveDescription(globalId, localInfo)
	if err != nil {
		log.Error("ERROR while saving localized description: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(localInfo)
}

// UpdateDescription is the handler for PUT /organizations/{globalid}/description
// Updates the description for this organization for this langKey
func (api OrganizationsAPI) UpdateDescription(w http.ResponseWriter, r *http.Request) {
	globalId := mux.Vars(r)["globalid"]
	var localInfo organization.LocalizedInfoText

	if err := json.NewDecoder(r.Body).Decode(&localInfo); err != nil {
		log.Debug("Error decoding the localized description:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	langKey := parseLangKey(localInfo.LangKey)
	if langKey == "" {
		log.Debug("Error decoding the localized description: invalid language key")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	localInfo.LangKey = langKey
	descriptionMgr := organization.NewDescriptionManager(r)
	err := descriptionMgr.UpdateDescription(globalId, localInfo)
	if err != nil {
		log.Error("ERROR while updating localized description: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(localInfo)
}

// AcceptOrganizationInvite is the handler for POST /organizations/{globalid}/organizations/{invitingorg}/roles/{role}
// Accept the organization invite for one of your organizations
func (api OrganizationsAPI) AcceptOrganizationInvite(w http.ResponseWriter, r *http.Request) {
	invitedorgname := mux.Vars(r)["globalid"]
	role := mux.Vars(r)["role"]
	invitingorgname := mux.Vars(r)["invitingorg"]

	var j invitations.JoinOrganizationInvitation

	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgReqMgr := invitations.NewInvitationManager(r)
	orgMgr := organization.NewManager(r)

	if !orgMgr.Exists(invitedorgname) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	invitingorg, err := orgMgr.GetByName(invitingorgname)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	orgRequest, err := orgReqMgr.Get(invitedorgname, invitingorgname, role, invitations.RequestPending)
	if err != nil {
		log.Error("error while trying to get invitation for organization")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if orgRequest == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if orgRequest.Role == invitations.RoleOrgOwner {
		if err := orgMgr.SaveOrgOwner(invitingorg, invitedorgname); err != nil {
			log.Error("Failed to save organization owner: ", invitedorgname)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		if err := orgMgr.SaveOrgMember(invitingorg, invitedorgname); err != nil {
			log.Error("Failed to save organization member: ", invitedorgname)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	orgRequest.Status = invitations.RequestAccepted

	if err := orgReqMgr.Save(orgRequest); err != nil {
		log.Error("Failed to update org request status: ", orgRequest.Organization)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(orgRequest)
}

// RejectOrganizationInvite is the handler for DELETE /organizations/{globalid}/organizations/{invitingorg}/role/{role}
// Reject the organization invite for one of your organizations
func (api OrganizationsAPI) RejectOrganizationInvite(w http.ResponseWriter, r *http.Request) {
	invitedorgname := mux.Vars(r)["globalid"]
	role := mux.Vars(r)["role"]
	invitingorgname := mux.Vars(r)["invitingorg"]

	orgReqMgr := invitations.NewInvitationManager(r)

	orgRequest, err := orgReqMgr.Get(invitedorgname, invitingorgname, role, invitations.RequestPending)
	if err != nil {
		log.Error("error while trying to load the invite")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if orgRequest == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	orgRequest.Status = invitations.RequestRejected

	if err := orgReqMgr.Save(orgRequest); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api OrganizationsAPI) inviteOrganization(w http.ResponseWriter, r *http.Request, role string, searchString string) {
	globalID := mux.Vars(r)["globalid"]

	// An organization can't invite itself.
	if searchString == globalID {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	orgMgr := organization.NewManager(r)
	org, err := orgMgr.GetByName(globalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			writeErrorResponse(w, http.StatusNotFound, "organization_not_found")
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}
	invitedOrg, err := orgMgr.GetByName(searchString)
	if err != nil {
		if err == mgo.ErrNotFound {
			writeErrorResponse(w, http.StatusNotFound, "invited_organization_not_found")
		} else {
			handleServerError(w, "getting invited organization", err)
		}
		return
	}

	if role == invitations.RoleOrgMember {
		for _, orgmembername := range org.OrgMembers {
			if orgmembername == invitedOrg.Globalid {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}
		}
	}
	for _, orgmemberName := range org.OrgOwners {
		if orgmemberName == invitedOrg.Globalid {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	// Create JoinRequest
	invitationMgr := invitations.NewInvitationManager(r)
	count, err := invitationMgr.CountByOrganization(globalID)
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if count >= maximumNumberOfInvitationsPerOrganization {
		log.Error("Reached invitation limit for organization ", globalID)
		writeErrorResponse(w, 422, "max_amount_of_invitations_reached")
		return
	}

	orgReq := &invitations.JoinOrganizationInvitation{
		Role:           role,
		Organization:   globalID,
		User:           invitedOrg.Globalid,
		Status:         invitations.RequestPending,
		Created:        db.DateTime(time.Now()),
		Method:         invitations.MethodWebsite,
		EmailAddress:   "",
		PhoneNumber:    "",
		Code:           "",
		IsOrganization: true,
	}

	if err = invitationMgr.Save(orgReq); err != nil {
		log.Error("Error inviting organization: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orgReq)
}

// AddIncludeSubOrgsOf is the handler for POST /organization/{globalid}/orgmembers/includesuborgs
// Include the suborganizations of the given organization in the member/owner hierarchy of this organization
func (api OrganizationsAPI) AddIncludeSubOrgsOf(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]

	includeSubOrgOf := struct {
		GlobalID string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&includeSubOrgOf); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := organization.NewManager(r)
	org, err := orgMgr.GetByName(globalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			writeErrorResponse(w, http.StatusNotFound, "organization_not_found")
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}

	if !orgMgr.Exists(includeSubOrgOf.GlobalID) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// check if the organization to add is already in the list
	for _, orgMembers := range org.IncludeSubOrgsOf {
		if orgMembers == includeSubOrgOf.GlobalID {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	// add the organization to the list
	err = orgMgr.AddIncludeSubOrgOf(org.Globalid, includeSubOrgOf.GlobalID)
	if handleServerError(w, "adding organization to 'includesuborgsof' list", err) {
		return
	}

	org, err = orgMgr.GetByName(globalID)
	if handleServerError(w, "getting organization", err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(org)
}

// RemoveIncludeSubOrgsOf is the handler for DELETE /organization/{globalid}/orgmembers/includesuborgs/{orgmember}
// Removes the suborganizations of the given organization from the member/owner hierarchy of this organization
func (api OrganizationsAPI) RemoveIncludeSubOrgsOf(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]
	orgMember := mux.Vars(r)["orgmember"]

	orgMgr := organization.NewManager(r)
	org, err := orgMgr.GetByName(globalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			writeErrorResponse(w, http.StatusNotFound, "organization_not_found")
		} else {
			handleServerError(w, "getting organization", err)
		}
		return
	}

	if !orgMgr.Exists(orgMember) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err = orgMgr.RemoveIncludeSubOrgOf(org.Globalid, orgMember)
	if handleServerError(w, "removing organization from 'includesuborgsof' list", err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UserIsMember is the handler for GET /organization/{globalid}/users/ismember/{username}
// Checks if the user has membership rights on the organization
func (api OrganizationsAPI) UserIsMember(w http.ResponseWriter, r *http.Request) {
	globalID := mux.Vars(r)["globalid"]
	username := mux.Vars(r)["username"]

	var isMember bool

	user, err := SearchUser(r, username)
	if err == mgo.ErrNotFound {
		user = nil
	} else if handleServerError(w, "getting user from database", err) {
		return
	}

	if user != nil {
		orgMgr := organization.NewManager(r)
		isMember, err = orgMgr.IsMember(globalID, username)
		if handleServerError(w, "checking if user is a member of the organization", err) {
			return
		}

		if !isMember {
			isMember, err = orgMgr.IsOwner(globalID, username)
			if handleServerError(w, "checking if user is an owner of the organization", err) {
				return
			}
		}
	}

	response := struct {
		IsMember bool
	}{
		IsMember: isMember,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

func writeErrorResponse(responseWriter http.ResponseWriter, httpStatusCode int, message string) {
	log.Debug(httpStatusCode, message)
	errorResponse := struct {
		Error string `json:"error"`
	}{Error: message}
	responseWriter.WriteHeader(httpStatusCode)
	json.NewEncoder(responseWriter).Encode(&errorResponse)
}

func handleServerError(responseWriter http.ResponseWriter, actionText string, err error) bool {
	if err != nil {
		log.Error("organizations_api: error while "+actionText, " - ", err)
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return true
	}
	return false
}

func SearchUser(r *http.Request, searchString string) (usr *user.User, err1 error) {
	userMgr := user.NewManager(r)
	usr, err1 = userMgr.GetByName(searchString)
	if err1 == mgo.ErrNotFound {
		valMgr := validationdb.NewManager(r)
		validatedPhonenumber, err2 := valMgr.GetByPhoneNumber(searchString)
		if err2 == mgo.ErrNotFound {
			validatedEmailAddress, err3 := valMgr.GetByEmailAddress(searchString)
			if err3 != nil {
				return nil, err3
			}
			return userMgr.GetByName(validatedEmailAddress.Username)
		}
		return userMgr.GetByName(validatedPhonenumber.Username)
	}
	return usr, err1
}

// parseLangKey return the first 2 characters of a string in lowercase. If the string is empty or has only 1 character, and empty string is returned
func parseLangKey(rawKey string) string {
	if len(rawKey) < 2 {
		return ""
	}
	return strings.ToLower(string(rawKey[0:2]))
}
