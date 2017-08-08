package organization

type OrganizationUser struct {
	User          string     `json:"user"`
	Role          string     `json:"role"`
	MissingScopes []string   `json:"missingscopes"`
}
