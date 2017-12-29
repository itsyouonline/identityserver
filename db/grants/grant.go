package grants

// Grant is a custom, application defined tag.
type Grant string

// SavedGrants links a username and globalid to stored grants
type SavedGrants struct {
	Username string  `json:"username"`
	GlobalID string  `json:"globalid"`
	Grants   []Grant `json:"grants"`
}
