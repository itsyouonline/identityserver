package grants

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// grantPrefix is the prefix added to grants when they are stored in access tokens / jwts.
	grantPrefix = "grant:"
	// grantMaxLen is the maximum length of a grant, in bytes.
	grantMaxLen = 100
)

var (
	// ErrGrantTooLarge indicates that the size of the grant in bytes is larger than grantMaxLen
	ErrGrantTooLarge = fmt.Errorf("Grant size too large, max %d bytes", grantMaxLen)
)

// Grant is a custom, application defined tag.
type Grant string

// SavedGrants links a username and globalid to stored grants
type SavedGrants struct {
	Username string  `json:"username"`
	GlobalID string  `json:"globalid"`
	Grants   []Grant `json:"grants"`
}

// Validate validates a raw grant. Validation is successfull if the error is nil
func (grant *Grant) Validate() error {
	g := string(*grant)
	if len(g) > grantMaxLen {
		return ErrGrantTooLarge
	}
	if len(g) < 2 {
		return errors.New("grants must have a minimum size of 2 bytes")
	}
	for i := 0; i < len(g); i++ {
		// First Check if a byte represents an ascii number
		if g[i] >= 48 && g[i] <= 57 {
			continue
		}
		// Now check lowercase letters
		if g[i] >= 97 && g[i] <= 122 {
			continue
		}
		// And uppercase letters
		if g[i] >= 65 && g[i] <= 90 {
			continue
		}
		// Also allow dashes and underscores
		if g[i] == "_"[0] || g[i] == "-"[0] {
			continue
		}
		return errors.New("invalid character in grant")
	}
	// Seems we are all good
	return nil
}

// FullName ensures that a grant starts with the grantPrefix
func FullName(grant Grant) string {
	if strings.HasPrefix(string(grant), grantPrefix) {
		return string(grant)
	}
	return grantPrefix + string(grant)
}
