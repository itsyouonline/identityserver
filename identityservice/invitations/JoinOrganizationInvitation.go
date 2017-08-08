package invitations

import (
	"reflect"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/db/validation"
	"gopkg.in/mgo.v2/bson"
)

type InvitationStatus string

const (
	RequestPending  InvitationStatus = "pending"
	RequestAccepted InvitationStatus = "accepted"
	RequestRejected InvitationStatus = "rejected"
)

const (
	RoleMember    = "member"
	RoleOwner     = "owner"
	RoleOrgMember = "orgmember"
	RoleOrgOwner  = "orgowner"
)

type InviteMethod string

const (
	MethodWebsite InviteMethod = "website"
	MethodEmail   InviteMethod = "email"
	MethodPhone   InviteMethod = "phone"
)
const (
	InviteUrl = "https://%s/login#/organizationinvite/%s"
)

//JoinOrganizationInvitation defines an invitation to join an organization
type JoinOrganizationInvitation struct {
	ID             bson.ObjectId    `json:"-" bson:"_id,omitempty"`
	Organization   string           `json:"organization"`
	Role           string           `json:"role"`
	User           string           `json:"user"`
	Status         InvitationStatus `json:"status"`
	Created        db.DateTime      `json:"created"`
	Method         InviteMethod     `json:"method"`
	EmailAddress   string           `json:"emailaddress"`
	PhoneNumber    string           `json:"phonenumber"`
	Code           string           `json:"-"`
	IsOrganization bool             `json:"isorganization"`
}

func ParseInvitationType(invitationType string) string {
	val := reflect.ValueOf(RequestAccepted).String()
	if val == invitationType {
		return val
	}
	val = reflect.ValueOf(RequestRejected).String()
	if val == invitationType {
		return val
	}
	return reflect.ValueOf(RequestPending).String()
}

type JoinOrganizationInvitationView struct {
	Organization   string           `json:"organization"`
	Role           string           `json:"role"`
	User           string           `json:"user"`
	Status         InvitationStatus `json:"status"`
	Created        db.DateTime      `json:"created"`
	Method         InviteMethod     `json:"method"`
	EmailAddress   string           `json:"emailaddress"`
	PhoneNumber    string           `json:"phonenumber"`
	IsOrganization bool             `json:"isorganization"`
}

func (inv *JoinOrganizationInvitation) ConvertToView(usrMgr *user.Manager, valMgr *validation.Manager) (*JoinOrganizationInvitationView, error) {
	vw := &JoinOrganizationInvitationView{}
	vw.Organization = inv.Organization
	vw.Role = inv.Role
	vw.Status = inv.Status
	vw.Created = inv.Created
	vw.Method = inv.Method
	vw.EmailAddress = inv.EmailAddress
	vw.PhoneNumber = inv.PhoneNumber
	vw.IsOrganization = inv.IsOrganization

	var err error
	vw.User, err = organization.ConvertUsernameToIdentifier(inv.User, usrMgr, valMgr)
	// user can be empty if invited through email or phone number
	if db.IsNotFound(err) {
		err = nil
	}
	return vw, err
}
