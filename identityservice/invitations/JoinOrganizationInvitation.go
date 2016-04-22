package invitations

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type InvitationStatus string

const (
	RequestPending  InvitationStatus = "pending"
	RequestAccepted InvitationStatus = "accepted"
	RequestRejected InvitationStatus = "rejected"
)

const (
	RoleMember = "member"
	RoleOwner  = "owner"
)

//JoinOrganizationInvitation defines an invitation to join an organization
type JoinOrganizationInvitation struct {
	ID           bson.ObjectId    `json:"-" bson:"_id,omitempty"`
	Organization string           `json:"organization"`
	Role         string           `json:"role"`
	User         string           `json:"user"`
	Status       InvitationStatus `json:"status"`
	Created      time.Time        `json:"created"`
}
