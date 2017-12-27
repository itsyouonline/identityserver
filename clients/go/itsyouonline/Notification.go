package itsyouonline

import (
	"gopkg.in/validator.v2"
)

type Notification struct {
	Approvals        []JoinOrganizationInvitation `json:"approvals" validate:"nonzero"`
	ContractRequests []ContractSigningRequest     `json:"contractRequests" validate:"nonzero"`
	Invitations      []JoinOrganizationInvitation `json:"invitations" validate:"nonzero"`
	Missingscopes    []MissingScopes              `json:"missingscopes" validate:"nonzero"`
}

func (s Notification) Validate() error {

	return validator.Validate(s)
}
