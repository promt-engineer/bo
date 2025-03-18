package requests

type Pagination struct {
	Filters map[string]interface{} `json:"filters,omitempty" form:"filters,omitempty"`
	Order   string                 `json:"order,omitempty" form:"order,omitempty"`
	Limit   int                    `json:"limit" form:"limit"`
	Offset  int                    `json:"offset" form:"offset"`
	Total   int64                  `json:"total" form:"total"`
}
