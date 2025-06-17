package response

import "time"

type AccessTokenResponse struct {
	Issuer       string    `json:"issuer"`
	Audience     string    `json:"audience"`
	Subject      string    `json:"subject"`
	ExpirationAt time.Time `json:"expiration_at"`
	IssuedAt     time.Time `json:"issued_at"`
	Token        string    `json:"token"`
}
