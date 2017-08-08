package organization

import (
	"regexp"

	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/db/validation"

	"gopkg.in/validator.v2"
)

type Organization struct {
	DNS              []string        `json:"dns"`
	Globalid         string          `json:"globalid"`
	Members          []string        `json:"members"`
	Owners           []string        `json:"owners"`
	PublicKeys       []string        `json:"publicKeys"`
	SecondsValidity  int             `json:"secondsvalidity"`
	OrgOwners        []string        `json:"orgowners"`  //OrgOwners are other organizations that are owner of this organization
	OrgMembers       []string        `json:"orgmembers"` //OrgMembers are other organizations that are member of this organization
	RequiredScopes   []RequiredScope `json:"requiredscopes"`
	IncludeSubOrgsOf []string        `json:"includesuborgsof"`
}

// IsValid performs basic validation on the content of an organizations fields
func (org *Organization) IsValid() bool {
	regex, _ := regexp.Compile(`^[a-z\d\-_\s]{3,150}$`)
	return validator.Validate(org) == nil && regex.MatchString(org.Globalid)
}

func (org *Organization) IsValidSubOrganization() bool {
	regex, _ := regexp.Compile(`^[a-z\d\-_\s\.]{3,150}$`)
	return validator.Validate(org) == nil && regex.MatchString(org.Globalid)
}

func (org *Organization) ConvertToView(usrMgr *user.Manager, valMgr *validation.Manager) (*OrganizationView, error) {
	view := &OrganizationView{}
	view.DNS = org.DNS
	view.Globalid = org.Globalid
	view.PublicKeys = org.PublicKeys
	view.SecondsValidity = org.SecondsValidity
	view.OrgOwners = org.OrgOwners
	view.OrgMembers = org.OrgMembers
	view.RequiredScopes = org.RequiredScopes
	view.IncludeSubOrgsOf = org.IncludeSubOrgsOf

	var err error
	view.Members, err = ConvertUsernamesToIdentifiers(org.Members, valMgr)
	if err != nil {
		return view, err
	}
	view.Owners, err = ConvertUsernamesToIdentifiers(org.Owners, valMgr)

	return view, err
}

func ConvertUsernamesToIdentifiers(usernames []string, valMgr *validation.Manager) ([]string, error) {
	identifiers := []string{}
	checkedUsers := map[string]bool{}
	for _, u := range usernames {
		checkedUsers[u] = false
	}
	emails, err := valMgr.GetValidatedEmailAddressesByUsernames(usernames)
	if err != nil {
		return identifiers, err
	}
	for _, validatedEmail := range emails {
		if !checkedUsers[validatedEmail.Username] {
			identifiers = append(identifiers, validatedEmail.EmailAddress)
			checkedUsers[validatedEmail.Username] = true
		}
	}
	checkPhoneUsernames := []string{}
	for username, checked := range checkedUsers {
		if !checked {
			checkPhoneUsernames = append(checkPhoneUsernames, username)
		}
	}
	validatedPhoneNumbers, err := valMgr.GetValidatedPhoneNumbersByUsernames(checkPhoneUsernames)
	if err != nil {
		return identifiers, err
	}
	for _, validatedPhone := range validatedPhoneNumbers {
		if !checkedUsers[validatedPhone.Username] {
			identifiers = append(identifiers, validatedPhone.Phonenumber)
			checkedUsers[validatedPhone.Username] = true
		}
	}
	return identifiers, nil
}

// MapUsernamesToIdentifiers returns a map with as key the validated information (identifier) and as value the username
func MapUsernamesToIdentifiers(usernames []string, valMgr *validation.Manager) (map[string]string, error) {
	identifiers := map[string]string{}
	emails, err := valMgr.GetValidatedEmailAddressesByUsernames(usernames)
	if err != nil {
		return identifiers, err
	}
	for _, validatedEmail := range emails {
		identifiers[validatedEmail.EmailAddress] = validatedEmail.Username
	}
	validatedPhoneNumbers, err := valMgr.GetValidatedPhoneNumbersByUsernames(usernames)
	if err != nil {
		return identifiers, err
	}
	for _, validatedPhone := range validatedPhoneNumbers {
		identifiers[ validatedPhone.Phonenumber] = validatedPhone.Username
	}
	return identifiers, nil
}

func ConvertUsernameToIdentifier(username string, usrMgr *user.Manager, valMgr *validation.Manager) (string, error) {
	userIdentifier := username
	usr, err := usrMgr.GetByName(username)
	if err != nil {
		return userIdentifier, err
	}
	// check for a validated email address
	for _, email := range usr.EmailAddresses {
		validated, err := valMgr.IsEmailAddressValidated(username, email.EmailAddress)
		if err != nil {
			return userIdentifier, err
		}
		if validated {
			return email.EmailAddress, err
		}
	}
	// try the phone numbers
	for _, phone := range usr.Phonenumbers {
		validated, err := valMgr.IsPhonenumberValidated(username, phone.Phonenumber)
		if err != nil {
			return userIdentifier, err
		}
		if validated {
			return phone.Phonenumber, err
		}
	}
	// No verified email or phone number. Fallback to username
	return userIdentifier, err
}
func ConvertIdentifierToUsername(identifier string, valMgr *validation.Manager) (string, error) {
	email, err := valMgr.GetByEmailAddress(identifier)
	if err == nil {
		return email.Username, err
	} else if valMgr.IsErrNotFound(err) {
		phone, err := valMgr.GetByPhoneNumber(identifier)
		if err == nil {
			return phone.Username, err
		} else {
			return identifier, err
		}
	}
	return identifier, err
}

type OrganizationView struct {
	DNS              []string        `json:"dns"`
	Globalid         string          `json:"globalid"`
	Members          []string        `json:"members"`
	Owners           []string        `json:"owners"`
	PublicKeys       []string        `json:"publicKeys"`
	SecondsValidity  int             `json:"secondsvalidity"`
	OrgOwners        []string        `json:"orgowners"`  //OrgOwners are other organizations that are owner of this organization
	OrgMembers       []string        `json:"orgmembers"` //OrgMembers are other organizations that are member of this organization
	RequiredScopes   []RequiredScope `json:"requiredscopes"`
	IncludeSubOrgsOf []string        `json:"includesuborgsof"`
}
