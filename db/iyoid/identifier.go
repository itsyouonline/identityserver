package iyoid

// Identifier represents the mapping of an ID to a username, with the azp of the party
// which is authorized to see it
type Identifier struct {
	Username string   `json:"username"`
	IyoIDs   []string `json:"iyoids"`
	Azp      string   `json:"azp"` // Only this party can see this ID
}
