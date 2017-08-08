package organization

type OrganizationUser struct {
	Username      string     `json:"username"`
	Role          string     `json:"role"`
	MissingScopes []string   `json:"missingscopes"`
}
