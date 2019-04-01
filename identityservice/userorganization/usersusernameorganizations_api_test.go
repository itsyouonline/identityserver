package userorganization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterOrgs(t *testing.T) {
	type testcase struct {
		scopes   []string
		orgs     []string
		expected []string
	}
	testcases := []testcase{
		{orgs: []string{}, scopes: []string{}, expected: []string{}},
		{orgs: []string{"parentorg", "parentorg.suborg1", "parentorg.suborg2"}, scopes: []string{"user:organizations:parentorg"}, expected: []string{"parentorg", "parentorg.suborg1", "parentorg.suborg2"}},
		{orgs: []string{"parentorg", "parentorg.suborg1", "parentorg.suborg2"}, scopes: []string{"user:organizations:parent"}, expected: []string{}},
		{orgs: []string{"parentorg.suborg1", "parentorg.suborg2"}, scopes: []string{"user:organizations:parentorg"}, expected: []string{"parentorg.suborg1", "parentorg.suborg2"}},
		{orgs: []string{"parentorg.suborg1.child", "parentorg.suborg2.child"}, scopes: []string{"user:organizations:parentorg.suborg1", "user:organizations:parentorg.suborg2"}, expected: []string{"parentorg.suborg1.child", "parentorg.suborg2.child"}},
	}
	for _, test := range testcases {
		assert.Equal(t, test.expected, filterOrgs(test.orgs, test.scopes), "Failed")
	}
}
