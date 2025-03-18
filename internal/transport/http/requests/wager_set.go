package requests

type PaginateWagerSetRequest struct {
	Limit int `json:"limit" form:"limit" validate:"required"`
	Page  int `json:"page" form:"page" validate:"required"`
}

type CreateWagerSetRequest struct {
	Name         string  `json:"name" form:"name" validate:"required"`
	WagerLevels  []int64 `json:"wager_levels" form:"wager_levels" validate:"required"`
	DefaultWager int64   `json:"default_wager" form:"default_wager" validate:"required"`
}

type UpdateWagerSetRequest struct {
	Name         string  `json:"name" form:"name" validate:"required"`
	WagerLevels  []int64 `json:"wager_levels" form:"wager_levels" validate:"required"`
	DefaultWager int64   `json:"default_wager" form:"default_wager" validate:"required"`
	IsActive     *bool   `json:"is_active" form:"is_active" validate:"required"`
}
