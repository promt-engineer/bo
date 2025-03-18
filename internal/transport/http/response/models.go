package response

type HealthResponse struct {
	Success string `json:"success"`
}

type InfoResponse struct {
	Tag string `json:"tag"`
}
