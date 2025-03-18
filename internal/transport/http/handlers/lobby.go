package handlers

import (
	"backoffice/internal/services"
	"backoffice/internal/transport/http/requests"
	"backoffice/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

type lobbyHandler struct {
	lobbyService *services.LobbyService
}

func NewLobbyHTTPHandler(lobbyService *services.LobbyService) *lobbyHandler {
	return &lobbyHandler{
		lobbyService: lobbyService,
	}
}

func (h *lobbyHandler) Register(router *gin.RouterGroup) {
	lobby := router.Group("lobby")

	lobby.POST("start_game", h.startGame)
}

// @Summary Generate startGame lobby link.
// @Tags lobby
// @Consume application/json
// @Description Processes the request to the game lobby and returns the URL to access the startGame.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.LobbyRequest true "CreateLobbyRequest"
// @Success 200 {object} map[string]interface{}
// @Router /api/lobby/start_game [post].
func (h *lobbyHandler) startGame(ctx *gin.Context) {
	req := requests.LobbyRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)
	}

	link, err := h.lobbyService.Lobby(ctx, req)
	if err != nil {
		response.ServerError(ctx, err, nil)
		return
	}

	response.OK(ctx, link, nil)
}
