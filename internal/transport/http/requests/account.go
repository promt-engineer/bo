package requests

type CreateAccountRequest struct {
	ID         string `json:"id" validate:"required"`
	Token      string `json:"token" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name"`
	RoleID     string `json:"role_id" validate:"required"`
	OperatorID string `json:"operator_id"`
}

type UpdateAccountRequest struct {
	//ID        string `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	//RoleID    string `json:"role_id" validate:"required"`
}
