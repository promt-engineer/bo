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

type wagerSetHandler struct {
	wagerSetService *services.WagerSetService
}

func NewWagerSetHandler(wagerSetService *services.WagerSetService) *wagerSetHandler {
	return &wagerSetHandler{wagerSetService: wagerSetService}
}

func (h *wagerSetHandler) Register(router *gin.RouterGroup) {
	wagerSets := router.Group("wager_set")

	wagerSets.GET("", h.all)
	wagerSets.POST("", h.create)

	wagerSet := wagerSets.Group(":id")
	{
		wagerSet.GET("", h.get)
		wagerSet.PUT("", h.update)
		wagerSet.DELETE("", h.delete)
	}
}

// @Summary Get wager sets.
// @Tags wager_set
// @Consume application/json
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param limit query int true "rows limit"
// @Param page query int true "page"
// @Success 200 {object} response.Response{data=[]entities.WagerSet}
// @Router /api/wager_set [get].
func (h *wagerSetHandler) all(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.PaginateWagerSetRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	paginate, err := h.wagerSetService.Paginate(ctx, session.OrganizationID, map[string]interface{}{}, req.Limit, req.Page)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, paginate, nil)
}

// @Summary Create wager set.
// @Tags wager_set
// @Consume application/json
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param  data body requests.CreateWagerSetRequest true "requests.CreateWagerSetRequest"
// @Success 200  {object} response.Response{data=entities.WagerSet}
// @Router /api/wager_set [post].
func (h *wagerSetHandler) create(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.CreateWagerSetRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	ws, err := h.wagerSetService.Create(ctx, session.OrganizationID, req.Name, req.WagerLevels, req.DefaultWager)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, ws, nil)
}

// @Summary Update wager set.
// @Tags wager_set
// @Consume application/json
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "wager_set_id"
// @Param data body requests.UpdateWagerSetRequest true "requests.UpdateWagerSetRequest"
// @Success 200 {object} response.Response{data=entities.WagerSet}
// @Router /api/wager_set/{id} [put].
func (h *wagerSetHandler) update(ctx *gin.Context) {
	req := &requests.UpdateWagerSetRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	wsID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	ws, err := h.wagerSetService.Update(ctx, wsID, req.Name, req.WagerLevels, req.DefaultWager, *req.IsActive)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, ws, nil)
}

// @Summary Get wager set.
// @Tags wager_set
// @Consume application/json
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "wager_set_id"
// @Success 200 {object} response.Response{data=entities.WagerSet}
// @Router /api/wager_set/{id} [get].
func (h *wagerSetHandler) get(ctx *gin.Context) {
	wsID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	ws, err := h.wagerSetService.Get(ctx, wsID)
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

// @Summary Delete wager set.
// @Tags wager_set
// @Consume application/json
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "wager_set_id"
// @Success 204
// @Router /api/wager_set/{id} [delete].
func (h *wagerSetHandler) delete(ctx *gin.Context) {
	wsID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	err = h.wagerSetService.Delete(ctx, wsID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.NoContent(ctx)
}
