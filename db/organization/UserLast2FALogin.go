package organization

import	"github.com/itsyouonline/identityserver/db"

type UserLast2FALogin struct {
  Globalid    string       `json:"globalid"`
  Username    string       `json:"username"`
  Last2FA     db.DateTime  `json:"last2fa"`
}
