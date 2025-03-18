package entities

type Pagination[T any] struct {
	Items       []*T `json:"items"`
	CurrentPage int  `json:"current_page"`
	Limit       int  `json:"limit"`
	Total       int  `json:"total"`
}

func PaginationSubstituteItems[O, N any](p Pagination[O], items []*N) Pagination[N] {
	return Pagination[N]{
		Items:       items,
		CurrentPage: p.CurrentPage,
		Limit:       p.Limit,
		Total:       p.Total,
	}
}
