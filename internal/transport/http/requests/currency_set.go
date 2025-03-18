package requests

type PaginateCurrencySetRequest struct {
	Limit int `json:"limit" form:"limit" validate:"required"`
	Page  int `json:"page" form:"page" validate:"required"`
}

type CreateCurrencySetRequest struct {
	Name       string   `json:"name" form:"name" validate:"required"`
	Currencies []string `json:"currencies" form:"currencies" validate:"required"`
}

type UpdateCurrencySetRequest struct {
	Name       string   `json:"name" form:"name" validate:"required"`
	Currencies []string `json:"currencies" form:"currencies" validate:"required"`
	IsActive   *bool    `json:"is_active" form:"is_active" validate:"required"`
}
