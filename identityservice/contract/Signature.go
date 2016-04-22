package contract

import "time"

type Signature struct {
	Date      time.Time `json:"date"`
	PublicKey string    `json:"publicKey"`
	Signature string    `json:"signature"`
	SignedBy  string    `json:"signedBy"`
}
