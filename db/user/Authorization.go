package user

import "strings"

// Authorization defines what userinformation is authorized to be seen by an organization
// For an explanation about scopes and scopemapping, see https://github.com/itsyouonline/identityserver/blob/master/docs/oauth2/scopes.md
type Authorization struct {
	Addresses      []AuthorizationMap           `json:"addresses,omitempty"`
	BankAccounts   []AuthorizationMap           `json:"bankaccounts,omitempty"`
	DigitalWallet  []DigitalWalletAuthorization `json:"digitalwallet,omitempty"`
	EmailAddresses []AuthorizationMap           `json:"emailaddresses,omitempty"`
	Facebook       bool                         `json:"facebook,omitempty"`
	Github         bool                         `json:"github,omitempty"`
	GrantedTo      string                       `json:"grantedTo"`
	Organizations  []string                     `json:"organizations"`
	Phonenumbers   []AuthorizationMap           `json:"phonenumbers,omitempty"`
	PublicKeys     []AuthorizationMap           `json:"publicKeys,omitempty"`
	Username       string                       `json:"username"`
	Name           bool                         `json:"name"`
	OwnerOf        OwnerOf                      `json:"ownerof,omitempty"`
	Avatars        []AuthorizationMap           `json:"avatars,omitempty"`
	KeyStore       bool                         `json:"keystore,omitempty"`
}

type AuthorizationMap struct {
	RequestedLabel string `json:"requestedlabel"`
	RealLabel      string `json:"reallabel"`
	Scope          string `json:"scope" bson:"scope,omitempty"` // "write" or nothing (for now)
}

type DigitalWalletAuthorization struct {
	AuthorizationMap
	Currency string `json:"currency"`
}

type OwnerOf struct {
	EmailAddresses []string `json:"emailaddresses"`
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
			if authorization.ContainsOrganization(requestedorgid) {
				authorizedScopes = append(authorizedScopes, scope)
			}
		}
		if scope == "user:github" && authorization.Github {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if scope == "user:facebook" && authorization.Facebook {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:address") && LabelledPropertyIsAuthorized(scope, "user:address", authorization.Addresses) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:bankaccount") && LabelledPropertyIsAuthorized(scope, "user:bankaccount", authorization.BankAccounts) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:digitalwalletaddress") && DigitalWalletIsAuthorized(scope, "user:digitalwalletaddress", authorization.DigitalWallet) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:email") && LabelledPropertyIsAuthorized(scope, "user:email", authorization.EmailAddresses) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:phone") && LabelledPropertyIsAuthorized(scope, "user:phone", authorization.Phonenumbers) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:publickey") && LabelledPropertyIsAuthorized(scope, "user:publickey", authorization.PublicKeys) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:ownerof:email") && OwnerOfIsAuthorized(scope, "user:ownerof:email", authorization.OwnerOf.EmailAddresses) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if strings.HasPrefix(scope, "user:avatar") && LabelledPropertyIsAuthorized(scope, "user:avatar", authorization.Avatars) {
			authorizedScopes = append(authorizedScopes, scope)
		}
		if scope == "user:keystore" && authorization.KeyStore {
			authorizedScopes = append(authorizedScopes, scope)
		}
	}

	return
}

func (authorization Authorization) ContainsOrganization(globalid string) bool {
	for _, orgid := range authorization.Organizations {
		if orgid == globalid {
			return true
		}
	}
	return false
}

// LabelledPropertyIsAuthorized checks if a labelled property is authorized
func LabelledPropertyIsAuthorized(scope string, scopePrefix string, authorizedLabels []AuthorizationMap) (authorized bool) {
	if authorizedLabels == nil {
		return
	}
	if scope == scopePrefix {
		authorized = len(authorizedLabels) > 0
		return
	}
	if strings.HasPrefix(scope, scopePrefix+":") {
		split := strings.Split(strings.TrimPrefix(scope, scopePrefix+":"), ":")
		requestedLabel := split[0]
		requestedScope := split[len(split)-1]
		if requestedLabel == requestedScope {
			requestedScope = ""
		}
		for _, authorizationmap := range authorizedLabels {
			if (authorizationmap.RequestedLabel == requestedLabel ||
				authorizationmap.RequestedLabel == "main" && requestedLabel == "") &&
				(authorizationmap.Scope == requestedScope || requestedScope == "") {
				authorized = true
				return
			}
		}
	}
	return
}

// DigitalWalletIsAuthorized checks if a digital wallet is authorized
func DigitalWalletIsAuthorized(scope string, scopePrefix string, authorizedLabels []DigitalWalletAuthorization) (authorized bool) {
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

func OwnerOfIsAuthorized(scope string, scopePrefix string, authorizedOwnerOfs []string) bool {
	if authorizedOwnerOfs == nil {
		return false
	}
	for _, authorizedOwnerOf := range authorizedOwnerOfs {
		requestedOwnerOf := strings.TrimPrefix(scope, scopePrefix+":")
		if authorizedOwnerOf == requestedOwnerOf {
			return true
		}
	}
	return false
}
