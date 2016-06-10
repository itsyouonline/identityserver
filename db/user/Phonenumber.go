package user

import "regexp"

//Phonenumber defines a phonenumber and has functions for validation
type Phonenumber string

var (
	phoneRegex = regexp.MustCompile(`^\+[0-9]+$`)
)

//IsValid checks if a phonenumber is in a valid format
func (phonenumber Phonenumber) IsValid() (valid bool) {
	valid = true
	valid = valid && len(phonenumber) < 51
	return phoneRegex.Match([]byte(phonenumber))
}
