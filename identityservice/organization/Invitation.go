package organization

import (
	"time"
)

type Invitation struct {
	Created time.Time `json:"created"`
	Role    string    `json:"role"`
	User    string    `json:"user"`
}
