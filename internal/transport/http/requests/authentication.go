package requests

type AuthenticateRequest struct {
	ID       string `json:"id"`
	Token    string `json:"token"`
	Provider string `json:"provider"`
}

type ChangePasswordRequest struct {
	Password        string `json:"password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,eqfield=ConfirmPassword,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

type ResetPassword struct {
	NewPassword     string `json:"new_password" validate:"required,eqfield=ConfirmPassword,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}
