package requests

type UpsertPermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Subject     string `json:"subject" validate:"required"`
	Endpoint    string `json:"endpoint" validate:"required"`
	Action      string `json:"action" validate:"required"`
}
