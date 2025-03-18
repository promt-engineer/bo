package handlers

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/services"
	"backoffice/internal/transport/http/requests"
	"backoffice/internal/transport/http/response"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type currencySetHandler struct {
	currencySetService *services.CurrencySetService
}

func NewCurrencySetHandler(currencySetService *services.CurrencySetService) *currencySetHandler {
	return &currencySetHandler{
		currencySetService: currencySetService,
	}
}

func (h *currencySetHandler) Register(router *gin.RouterGroup) {
	currencySets := router.Group("currency_set")

	currencySets.GET("", h.all)
	currencySets.POST("", h.create)

	currencySet := currencySets.Group(":id")
	{
		currencySet.GET("", h.get)
		currencySet.PUT("", h.update)
		currencySet.DELETE("", h.delete)
	}
}

// DEPRECATE, DO NOT USE
func (h *currencySetHandler) all(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.PaginateCurrencySetRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	paginate, err := h.currencySetService.Paginate(ctx, session.OrganizationID, map[string]interface{}{}, req.Limit, req.Page)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, paginate, nil)
}

// DEPRECATE, DO NOT USE
func (h *currencySetHandler) create(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.CreateCurrencySetRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	ws, err := h.currencySetService.Create(ctx, session.OrganizationID, req.Name, req.Currencies)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, ws, nil)
}

// DEPRECATE, DO NOT USE
func (h *currencySetHandler) update(ctx *gin.Context) {
	req := &requests.UpdateCurrencySetRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	wsID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	ws, err := h.currencySetService.Update(ctx, wsID, req.Name, req.Currencies, *req.IsActive)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, ws, nil)
}

// DEPRECATE, DO NOT USE
func (h *currencySetHandler) get(ctx *gin.Context) {
	wsID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	ws, err := h.currencySetService.Get(ctx, wsID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, ws, nil)
}

// DEPRECATE, DO NOT USE
func (h *currencySetHandler) delete(ctx *gin.Context) {
	wsID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	err = h.currencySetService.Delete(ctx, wsID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.NoContent(ctx)
}
