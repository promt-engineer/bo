package response

type Pagination struct {
	Data    interface{}            `json:"data"`
	Filters map[string]interface{} `json:"filters,omitempty"`
	Limit   int                    `json:"limit,omitempty"`
	Offset  int                    `json:"offset,omitempty"`
	Total   int64                  `json:"total"`
}

func (p *Pagination) GetMeta() map[string]interface{} {
	mp := make(map[string]interface{})
	mp["filter"] = p.Filters
	mp["limit"] = p.Limit
	mp["offset"] = p.Offset
	mp["total"] = p.Total
	return mp
}
