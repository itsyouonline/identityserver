package user

import "strings"

// Authorization defines what userinformation is authorized to be seen by an organization
// For an explanation about scopes and scopemapping, see https://github.com/itsyouonline/identityserver/blob/master/docs/oauth2/scopes.md
type Authorization struct {
	Addresses      []AuthorizationMap `json:"addresses,omitempty"`
	BankAccounts   []AuthorizationMap `json:"bankaccounts,omitempty"`
	DigitalWallet  []AuthorizationMap `json:"digitalwallet,omitempty"`
	EmailAddresses []AuthorizationMap `json:"emailaddresses,omitempty"`
	Facebook       bool               `json:"facebook,omitempty"`
	Github         bool               `json:"github,omitempty"`
	GrantedTo      string             `json:"grantedTo"`
	Organizations  []string           `json:"organizations"`
	Phonenumbers   []AuthorizationMap `json:"phonenumbers,omitempty"`
	PublicKeys     []AuthorizationMap `json:"publicKeys,omitempty"`
	Username       string             `json:"username"`
	Name           bool               `json:"name"`
}

type AuthorizationMap struct {
	RequestedLabel string `json:"requestedlabel"`
	RealLabel      string `json:"reallabel"`
}

//FilterAuthorizedScopes filters the requested scopes to the ones this Authorization covers
func (authorization Authorization) FilterAuthorizedScopes(requestedscopes []string) (authorizedScopes []string) {
	authorizedScopes = make([]string, 0, len(requestedscopes))
	for _, rawscope := range requestedscopes {
		scope := strings.TrimSpace(rawscope)
		if scope == "user:name" && authorization.Name {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:memberof:") {
			requestedorgid := strings.TrimPrefix(scope, "user:memberof:")
			if authorization.containsOrganization(requestedorgid) {
				authorizedScopes = append(authorizedScopes, scope)
			}
		}
		if scope == "user:github" && authorization.Github {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if scope == "user:facebook" && authorization.Facebook {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:address") && labelledPropertyIsAuthorized(scope, "user:address", authorization.Addresses) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:bankaccount") && labelledPropertyIsAuthorized(scope, "user:bankaccount", authorization.BankAccounts) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:digitalwalletaddress") && labelledPropertyIsAuthorized(scope, "user:digitalwalletaddress", authorization.DigitalWallet) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:email") && labelledPropertyIsAuthorized(scope, "user:email", authorization.EmailAddresses) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:phone") && labelledPropertyIsAuthorized(scope, "user:phone", authorization.Phonenumbers) {
			authorizedScopes = append(authorizedScopes, scope)
		}
	}

	return
}

func (authorization Authorization) containsOrganization(globalid string) bool {
	for _, orgid := range authorization.Organizations {
		if orgid == globalid {
			return true
		}
	}
	return false
}

func labelledPropertyIsAuthorized(scope string, scopePrefix string, authorizedLabels []AuthorizationMap) (authorized bool) {
	if authorizedLabels == nil {
		return
	}
	if scope == scopePrefix {
		authorized = len(authorizedLabels) > 0
		return
	}
	if strings.HasPrefix(scope, scopePrefix+":") {
		requestedLabel := strings.Split(strings.TrimPrefix(scope, scopePrefix+":"), ":")[0]
		for _, authorizationmap := range authorizedLabels {
			if authorizationmap.RequestedLabel == requestedLabel || authorizationmap.RequestedLabel == "main" && requestedLabel == "" {
				authorized = true
				return
			}
		}
	}
	return
}
