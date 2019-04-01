package oauthservice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScopePresent(t *testing.T) {
	assert := assert.New(t)

	const (
		searchScope = "openid"
	)
	cases := []struct {
		scopes string
		isOIDC bool
	}{
		{"openid", true},
		{"s1,openid", true},
		{"s1,openid,s2", true},
		{"s1,s2", false},
		{"s1,open-id", false},
		{"", false},
	}

	for _, c := range cases {
		assert.Equal(c.isOIDC, scopePresent(c.scopes, searchScope), fmt.Sprintf("'%s' should return %t from isOIDC", c.scopes, c.isOIDC))
	}
}
