package organization

import "strings"

type Organization struct {
	DNS        		 	[]string `json:"dns"`
	Globalid   		 	string   `json:"globalid"`
	Members    		 	[]string `json:"members"`
	Owners     		 	[]string `json:"owners"`
	PublicKeys 		 	[]string `json:"publicKeys"`
	SecondsValidity	int			 `json:"secondsvalidity"`
}

// IsValid performs basic validation on the content of an organizations fields
//TODO: globalid should not contain ':,.'
func (c *Organization) IsValid() (valid bool) {
	valid = true
	globalIDLength := len(c.Globalid)
	valid = valid && (globalIDLength >= 3) && (globalIDLength <= 150) && c.Globalid == strings.ToLower(c.Globalid)
	return
}
