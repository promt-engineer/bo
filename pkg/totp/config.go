package totp

import "time"

type Config struct {
	Issuer     string
	SecretSize int
	Algorithm  Algorithm
	Period     time.Duration
	Digits     Digits
}
