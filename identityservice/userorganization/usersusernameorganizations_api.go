package userorganization

//The reason this api is not in the user package is because this would cause circular imports

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	organizationdb "github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/identityservice/invitations"
)

type UsersusernameorganizationsAPI struct {
}

func exists(value string, list []string) bool {
	for _, val := range list {
		if val == value {
			return true
		}
	}

	return false
}

// Get the list organizations a user is owner or member of
// It is handler for GET /users/{username}/organizations
func (api UsersusernameorganizationsAPI) Get(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	orgMgr := organizationdb.NewManager(r)

	orgs, err := orgMgr.AllByUserChain(username)
	if err != nil {
		log.Error("Failed to load all organizations the user belongs to: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type UserOrganizations struct {
		Member []string `json:"member"`
		Owner  []string `json:"owner"`
	}
	userOrgs := UserOrganizations{
		Member: []string{},
		Owner:  []string{},
	}

	userOrgs.Member = orgs
	userOrgs.Owner, userOrgs.Member, err = orgMgr.SplitOwnedOrgs(userOrgs.Member, username)
	if err != nil {
		log.Error("Failed to sort organizations to owner and member: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(&userOrgs)
}

// Accept membership in organization
// It is handler for POST /users/{username}/organizations/{globalid}/roles/{role}
func (api UsersusernameorganizationsAPI) globalidrolesrolePost(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	role := mux.Vars(r)["role"]
	organization := mux.Vars(r)["globalid"]

	var j invitations.JoinOrganizationInvitation

	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgReqMgr := invitations.NewInvitationManager(r)
	valMgr := validation.NewManager(r)

	var orgRequest *invitations.JoinOrganizationInvitation
	var err error

	switch j.Method {
	case invitations.MethodEmail:
		orgRequest, err = orgReqMgr.GetWithEmail(j.EmailAddress, organization, role, invitations.RequestPending)
		if err != nil {
			log.Error("error while trying to get invitation for email address: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if orgRequest == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		validatedEmail, e := valMgr.GetByEmailAddressValidatedEmailAddress(j.EmailAddress)
		if e != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		username = validatedEmail.Username
		break
	case invitations.MethodPhone:
		orgRequest, err = orgReqMgr.GetWithPhonenumber(j.PhoneNumber, organization, role, invitations.RequestPending)
		if err != nil {
			log.Error("error while trying to get invitation for phone number: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if orgRequest == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		validatedPhonenumber, e := valMgr.GetByPhoneNumber(j.PhoneNumber)
		if e != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		username = validatedPhonenumber.Username
		break
	default:
		orgRequest, err = orgReqMgr.Get(username, organization, role, invitations.RequestPending)
		if err != nil {
			log.Error("error while trying to get invitation for username: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if orgRequest == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	}

	// TODO: Save member
	orgMgr := organizationdb.NewManager(r)

	if org, err := orgMgr.GetByName(organization); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else {
		if invitations.RoleOwner == orgRequest.Role {
			// Accepted Owner role
			if err := orgMgr.SaveOwner(org, username); err != nil {
				log.Error("Failed to save owner: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else {
			// Accepted member role
			if err := orgMgr.SaveMember(org, username); err != nil {
				log.Error("Failed to save member: ", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
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

// Reject membership invitation in an organization.

// It is handler for DELETE /users/{username}/organizations/{globalid}/roles/{role}
func (api UsersusernameorganizationsAPI) globalidrolesroleDelete(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	role := mux.Vars(r)["role"]
	organization := mux.Vars(r)["globalid"]

	orgReqMgr := invitations.NewInvitationManager(r)

	orgRequest, err := orgReqMgr.Get(username, organization, role, invitations.RequestPending)
	if err != nil {
		log.Error("error while trying to load the invite")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// User was invited through phone or email
	// TODO: maybe think about a better way to handle this in the future.
	if orgRequest == nil {
		valMgr := validation.NewManager(r)
		validatedPhonenumbers, err := valMgr.GetByUsernameValidatedPhonenumbers(username)
		if err != nil {
			log.Error("error while loading verified phone numbers for user ", username)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		for _, number := range validatedPhonenumbers {
			orgRequest, err = orgReqMgr.GetWithPhonenumber(number.Phonenumber, organization, role, invitations.RequestPending)
			if err != nil {
				log.Error("error while trying to get invite for phone number")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if orgRequest != nil {
				break
			}
		}

		if orgRequest == nil {
			validatedEmailaddresses, err := valMgr.GetByUsernameValidatedEmailAddress(username)
			if err != nil {
				log.Error("error while loading verified phone numbers for user ", username)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			for _, email := range validatedEmailaddresses {
				orgRequest, err = orgReqMgr.GetWithEmail(email.EmailAddress, organization, role, invitations.RequestPending)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				if orgRequest != nil {
					break
				}
			}
		}
	}
	// All possible invite options exhausted, if we still haven't found the invite it just isn't there
	// This is possible when a user validates an email/phonennumber AFTER it was invited, then deletes
	// the verified email/phone and finally tries to decline the invitation before reloading the page
	if orgRequest == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	orgMgr := organizationdb.NewManager(r)

	if org, err := orgMgr.GetByName(organization); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else {
		if invitations.RoleOwner == orgRequest.Role {
			// Rejected Owner role
			if err := orgMgr.RemoveOwner(org, username); err != nil {
				log.Error("Failed to remove owner: ", username)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else {
			// Rejected member role
			if err := orgMgr.RemoveMember(org, username); err != nil {
				log.Error("Failed to reject member: ", username)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}

	orgRequest.Status = invitations.RequestRejected

	if err := orgReqMgr.Save(orgRequest); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
