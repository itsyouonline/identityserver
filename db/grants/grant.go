package grants

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// grantPrefix is the prefix added to grants when they are stored in access tokens / jwts.
	grantPrefix = "grant:"
	// grantMaxLen is the maximum length of a grant, in bytes.
	grantMaxLen = 100
	// grantMinLen is the minimum length of a grant, in bytes.
	grantMinLen = 2
)

var (
	// ErrGrantTooLarge indicates that the size of the grant in bytes is larger than grantMaxLen
	ErrGrantTooLarge = fmt.Errorf("Grant size too large, max %d bytes", grantMaxLen)
	// ErrGrantTooSmall indicates that the size of the grant in bytes is less than grantMinLen
	ErrGrantTooSmall = fmt.Errorf("Grant size too small, max %d bytes", grantMinLen)
	// ErrInvalidGrantCharacter indicates that a grant contains an invalid character
	ErrInvalidGrantCharacter = errors.New("Invalid character in grant")
	grantRegex               = regexp.MustCompile(`^[a-zA-Z0-9\-_\.]+$`)
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
	if len(g) < grantMinLen {
		return ErrGrantTooSmall
	}

	if grantRegex.MatchString(g) {
		return nil
	}
	return ErrInvalidGrantCharacter
}

// FullName ensures that a grant starts with the grantPrefix
func FullName(grant Grant) string {
	if strings.HasPrefix(string(grant), grantPrefix) {
		return string(grant)
	}
	return grantPrefix + string(grant)
}
