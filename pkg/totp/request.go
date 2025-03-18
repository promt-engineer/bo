package totp

type Request struct {
	TOTP string `json:"totp" validate:"required"`
}
