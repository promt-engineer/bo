package entities

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID `json:"id"`
	Account        *Account  `json:"account"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Currency       string    `json:"currency"`
}

func (s *Session) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}

func (s *Session) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s *Session) CleanUP() *Session {
	s.Account.TOTPSecret = ""
	s.Account.TOTPURL = ""

	return s
}
