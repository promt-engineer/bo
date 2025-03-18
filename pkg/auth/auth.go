package auth

import "time"

type Auth struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (t *Auth) Expired() bool {
	return t.ExpiredAt.Unix() < time.Now().Unix()
}
