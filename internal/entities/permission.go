package entities

import (
	"backoffice/internal/constants"
	"github.com/google/uuid"
	"net/http"
)

var actions = map[string]string{
	http.MethodGet:    constants.ActionView,
	http.MethodPut:    constants.ActionEdit,
	http.MethodPost:   constants.ActionCreate,
	http.MethodDelete: constants.ActionDelete,
}

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	Endpoint    string    `json:"endpoint"`
	Action      string    `json:"action"`
}

func (p *Permission) IsActionMatched(method string) bool {
	return p.Action == actions[method]
}
