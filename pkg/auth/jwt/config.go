package jwt

import "time"

type Config struct {
	HeaderName           string
	QueryName            string
	HeaderScheme         string
	Fingerprint          string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
	Issuer               string
}
