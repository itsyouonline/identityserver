package identifier

// Identifier represents the mapping of an ID to a username, with the azp of the party
// which is authorized to see it
type Identifier struct {
	Username string   `json:"username"`
	IDs      []string `json:"ids"`
	Azp      string   `json:"azp"` // Only this party can see this ID
}
