package organization

import (
	"regexp"
	"strings"

	"github.com/itsyouonline/identityserver/db/user"
)

type RequiredScope struct {
	Scope        string   `json:"scope"`
	AccessScopes []string `json:"accessscopes"`
}

var (
	accessScopeRegex = regexp.MustCompile(`organization:member|organization:owner`)
)

func (requiredScope RequiredScope) IsValid() bool {
	possibleScopes := []string{
		"user:name",
		"user:memberof:",
		"user:github",
		"user:facebook",
		"user:address",
		"user:bankaccount",
		"user:digitalwalletaddress",
		"user:email",
		"user:validated:email",
		"user:phone",
		"user:validated:phone",
		"user:publickey",
		"user:avatar",
		"user:keystore",
		"user:see",
	}
	valid := false
	for _, scope := range possibleScopes {
		if strings.HasPrefix(requiredScope.Scope, scope) {
			valid = true
		}
	}
	valid2 := false
	for _, accessScope := range requiredScope.AccessScopes {
		if accessScopeRegex.Match([]byte(accessScope)) {
			valid2 = true
		}
	}
	return valid && valid2
}

func (requiredScope *RequiredScope) IsAuthorized(authorization user.Authorization) bool {
	scope := requiredScope.Scope
	if scope == "user:name" && !authorization.Name {
		return false
	}
	if strings.HasPrefix(scope, "user:memberof:") {
		requestedOrganizationId := strings.TrimPrefix(scope, "user:memberof:")
		if !authorization.ContainsOrganization(requestedOrganizationId) {
			return false
		}
	}
	if scope == "user:github" && !authorization.Github {
		return false
	}
	if scope == "user:facebook" && !authorization.Facebook {
		return false
	}
	if strings.HasPrefix(scope, "user:address") && !user.LabelledPropertyIsAuthorized(scope, "user:address", authorization.Addresses) {
		return false
	}
	if strings.HasPrefix(scope, "user:bankaccount") && !user.LabelledPropertyIsAuthorized(scope, "user:bankaccount", authorization.BankAccounts) {
		return false
	}
	if strings.HasPrefix(scope, "user:digitalwalletaddress") && !user.DigitalWalletIsAuthorized(scope, "user:digitalwalletaddress", authorization.DigitalWallet) {
		return false
	}
	if strings.HasPrefix(scope, "user:email") && !user.LabelledPropertyIsAuthorized(scope, "user:email", authorization.EmailAddresses) {
		return false
	}
	if strings.HasPrefix(scope, "user:validated:email") && !user.LabelledPropertyIsAuthorized(scope, "user:validated:email", authorization.ValidatedEmailAddresses) {
		return false
	}
	if strings.HasPrefix(scope, "user:phone") && !user.LabelledPropertyIsAuthorized(scope, "user:phone", authorization.Phonenumbers) {
		return false
	}
	if strings.HasPrefix(scope, "user:validated:phone") && !user.LabelledPropertyIsAuthorized(scope, "user:validated:phone", authorization.Phonenumbers) {
		return false
	}
	if strings.HasPrefix(scope, "user:publickey") && !user.LabelledPropertyIsAuthorized(scope, "user:publickey", authorization.PublicKeys) {
		return false
	}
	if scope == "user:keystore" && !authorization.KeyStore {
		return false
	}
	if scope == "user:see" && !authorization.See {
		return false
	}
	return true
}
