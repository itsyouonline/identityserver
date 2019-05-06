package user

// TwoFactorAuthSettings is a Two Factor Authentication user Settings
type TwoFactorAuthSettings struct {
	Enable                bool `json:"enable"`
	SkipForOrgsIfOptional bool `json:"skipfororgsifoptional"`
}
