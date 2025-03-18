package dto

type Msg struct {
	Type    string `json:"type"`
	Payload interface{}
}
