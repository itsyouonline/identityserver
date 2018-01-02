package grants

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrantValidation(t *testing.T) {
	type testCase struct {
		rawGrant Grant
		valid    bool
	}
	testCases := []testCase{
		{"adfsaf", true},
		{" 55dasf", false},
		{"dsaf978e", true},
		{"grant", true},
		{"grant:test", false},
		{"grant_with_underscore", true},
		{"grant-with-dash", true},
		{"grant with space", false},
		{"g", false},
		{"g/adfe5", false},
		{"傀dsf", false},
		{"mail@provider", false},
		{"Unicode䄵love", false},
	}

	for i, test := range testCases {
		assert.Equal(t, test.rawGrant.Validate() == nil, test.valid, fmt.Sprintf("Testcase %d failed (%v)", i+1, test.rawGrant))
	}
}
