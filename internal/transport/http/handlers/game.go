package handlers

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/services"
	"backoffice/internal/transport/http/requests"
	"backoffice/internal/transport/http/response"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type gameHandler struct {
	gameService *services.GameService
	cfgSender   *services.ConfigSenderService
}

func NewGameHTTPHandler(gameService *services.GameService, cfgSender *services.ConfigSenderService) *gameHandler {
	return &gameHandler{
		gameService: gameService,
		cfgSender:   cfgSender,
	}
}

func (h *gameHandler) Register(route *gin.RouterGroup) {
	games := route.Group("game")
	{
		games.POST("", h.create)
		games.GET("", h.all)
		games.POST("search", h.search)
		game := games.Group(":id")
		{
			game.GET("", h.get)
			game.DELETE("", h.delete)
			game.PUT("", h.update)
		}
	}
}

// @Summary Get game list.
// @Tags game
// @Consume application/json
// @Description All games.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param limit query int true "rows limit"
// @Param offset query int true "rows offset"
// @Param order query string false "order field"
// @Success 200 {object} response.Response{data=[]entities.Game}
// @Router /api/game [get].
func (h *gameHandler) all(ctx *gin.Context) {
	req := &requests.Pagination{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	session := ctx.Value("session").(*entities.Session)
	games, total, err := h.gameService.Paginate(ctx, map[string]interface{}{"organization_id": session.OrganizationID}, req.Order, req.Limit, req.Offset)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	req.Total = total

	response.OK(ctx, games, req)
}

// @Summary Get games by filter.
// @Tags game
// @Consume application/json
// @Description Filtering can be done by organization_id or you can send an empty request body
// @Description and get all games
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param request body requests.GameListRequest true "Game get list request"
// @Success 200  {object} response.Response{data=[]entities.Game}
// @Router /api/game/search [post].
func (h *gameHandler) search(ctx *gin.Context) {
	req := &requests.GameListRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)
		return
	}

	filter := make(map[string]interface{})

	if req.OrganizationID != uuid.Nil {
		filter["organization_id"] = req.OrganizationID
	}

	game, err := h.gameService.GetGameList(ctx, filter)
	if err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, game, nil)
}

// @Summary Create a new game.
// @Tags game
// @Consume application/json
// @Description Create game.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param request body requests.GameRequest true "Game creation request"
// @Success 200 {object} response.Response{data=entities.Game}
// @Router /api/game [post]
func (h *gameHandler) create(ctx *gin.Context) {
	req := &requests.GameRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	game, err := h.gameService.Create(ctx, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, game, nil)
}

// @Summary Get game.
// @Tags game
// @Consume application/json
// @Description Get game information.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "id"
// @Success 200  {object} response.Response{data=entities.Game}
// @Router /api/game/{id} [get].
func (h *gameHandler) get(ctx *gin.Context) {
	gameID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	game, err := h.gameService.GetGame(ctx, gameID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)
			return
		}

		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, game, nil)
}

// @Summary Delete game.
// @Tags game
// @Consume application/json
// @Description delete game.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "id"
// @Success 204 "No Content"
// @Router /api/game/{id} [delete].
func (h *gameHandler) delete(ctx *gin.Context) {
	gameID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	err = h.gameService.Delete(ctx, gameID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.NoContent(ctx)
}

// @Summary Updates game information.
// @Tags game
// @Consume application/json
// @Description Updates game information.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "id"
// @Param request body requests.GameRequest true "Game update request"
// @Success 200  {object} response.Response{data=entities.Game}
// @Router /api/game/{id} [put].
func (h *gameHandler) update(ctx *gin.Context) {
	req := &requests.GameRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	gameID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	game, err := h.gameService.Update(ctx, gameID, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, game, nil)
}
