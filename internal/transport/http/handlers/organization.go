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
	"go.uber.org/zap"
)

type organizationHandler struct {
	organizationService *services.OrganizationService
	cfgSender           *services.ConfigSenderService
}

func NewOrganizationHandler(organizationService *services.OrganizationService, cfgSender *services.ConfigSenderService) *organizationHandler {
	return &organizationHandler{
		organizationService: organizationService,
		cfgSender:           cfgSender,
	}
}

func (h *organizationHandler) Register(route *gin.RouterGroup) {
	organizations := route.Group("organizations")
	{
		organizations.GET("", h.all)
		organizations.POST("", h.create)
		organizations.POST("get_organization_pair", h.getOrganizationPair)
		organizations.POST("add_organization_pair", h.addOrganizationPair)
		organizations.DELETE("delete_organization_pair", h.deleteOrganizationPair)
		organizations.POST("get_operator_pair", h.getOperatorPair)
		organizations.POST("add_operator_pair", h.addOperatorPair)
		organizations.DELETE("delete_operator_pair", h.deleteOperatorPair)

		organization := organizations.Group(":id")
		{
			organization.DELETE("", h.delete)
			organization.GET("", h.get)
			organization.GET("integrators", h.getIntegratorsByProvider)
			organization.GET("operators", h.getOperatorsByIntegrator)
			organization.PUT("", h.update)
			organization.GET("game", h.getGame)
			organization.POST("game", h.assignGames)
			organization.PUT("game", h.updateGame)
			organization.DELETE("game", h.revokeGames)
			organization.GET("wager_set", h.getGamesWagerSets)
			organization.POST("wager_set", h.addGameWagerSet)
			organization.PUT("wager_set", h.updateGameWagerSet)
			organization.DELETE("wager_set", h.deleteGameWagerSet)
		}
	}
}

// @Summary Get organizations list.
// @Tags organizations
// @Consume application/json
// @Description Available backoffice roles.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param limit query int true "rows limit"
// @Param offset query int true "rows offset"
// @Param order query string false "order field"
// @Success 200 {object} response.Response{data=[]entities.Organization}
// @Router /api/organizations [get].
func (h *organizationHandler) all(ctx *gin.Context) {
	req := &requests.Pagination{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	organizations, total, err := h.organizationService.Paginate(ctx, req.Filters, req.Order, req.Limit, req.Offset)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	req.Total = total

	response.OK(ctx, organizations, req)
}

// @Summary Get organization integrator list by provider.
// @Tags organizations
// @Consume application/json
// @Description Available backoffice roles.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Success 200 {object} response.Response{data=[]entities.Organization}
// @Router /api/organizations/{id}/integrators [get].
func (h *organizationHandler) getIntegratorsByProvider(ctx *gin.Context) {
	organizationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	organizations, err := h.organizationService.GetIntegratorsByProvider(ctx, organizationID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, organizations, nil)
}

// @Summary Get organization.
// @Tags organizations
// @Consume application/json
// @Description Get organization.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Success 200  {object} response.Response{data=entities.Organization}
// @Router /api/organizations/{id} [get].
func (h *organizationHandler) get(ctx *gin.Context) {
	organizationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	organization, err := h.organizationService.Get(ctx, organizationID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, organization, nil)
}

// @Summary Add new organization.
// @Tags organizations
// @Consume application/json
// @Description Create organization.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.UpsertOrganizationRequest true  "CreateOrganizationRequest"
// @Success 200 {object} response.Response{data=entities.Organization}
// @Router /api/organizations [post].
func (h *organizationHandler) create(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.UpsertOrganizationRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	organization, err := h.organizationService.Create(ctx, uint8(req.Status), req.Name, req.Type)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	if _, err = h.organizationService.Assign(ctx, session.Account.ID, organization.ID); err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, organization, nil)
}

// @Summary Update organization.
// @Tags organizations
// @Consume application/json
// @Description Update organization.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Param data body requests.UpsertOrganizationRequest true "UpdateOrganizationRequest"
// @Success 200 {object} response.Response{data=entities.Organization}
// @Router /api/organizations/{id} [put].
func (h *organizationHandler) update(ctx *gin.Context) {
	req := &requests.UpsertOrganizationRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	organizationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	organization, err := h.organizationService.Update(ctx, organizationID, uint8(req.Status), req.Name, req.Type)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, organization, nil)
}

// @Summary Delete organization.
// @Tags organizations
// @Consume application/json
// @Description Delete existing organization.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Success 204
// @Router /api/organizations/{id} [delete].
func (h *organizationHandler) delete(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	organizationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err = h.organizationService.Delete(ctx, session.Account.ID, organizationID)
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

// @Summary get integrator games.
// @Tags organizations
// @Consume application/json
// @Description get integrator games.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Success 200 {object} response.Response{data=[]entities.IntegratorGame}
// @Router /api/organizations/{id}/game [get].
func (h *organizationHandler) getGame(ctx *gin.Context) {
	integratorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	ig, err := h.organizationService.GetGames(ctx, integratorID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, ig, nil)
}

// @Summary add integrator game.
// @Tags organizations
// @Consume application/json
// @Description The provider_id field is required. If you send only the provider_id,
// @Description all games from the provider will be added. To add selected games, fill out the game_id.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Param data body requests.IntegratorGameRequest true "IntegratorGameRequest"
// @Success 200 {object} response.Response{data=[]entities.IntegratorGame}
// @Router /api/organizations/{id}/game [post].
func (h *organizationHandler) assignGames(ctx *gin.Context) {
	req := &requests.IntegratorGameRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	integratorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}
	zap.S().Info("init method cfgSender.Fire for integrator-games")
	ig, err := h.organizationService.AssignGames(ctx, integratorID, req.WagerSetID, req.GameID...)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, ig, nil)
}

// @Summary update integrator game.
// @Tags organizations
// @Consume application/json
// @Description The provider_id field is required. If you send only the provider_id,
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Param data body requests.UpdateIntegratorGameRequest true "IntegratorGameRequest"
// @Success 200 {object} response.Response{data=entities.IntegratorGame}
// @Router /api/organizations/{id}/game [put].
func (h *organizationHandler) updateGame(ctx *gin.Context) {
	req := &requests.UpdateIntegratorGameRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	integratorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}
	zap.S().Info("init method cfgSender.Fire for integrator-games")
	ig, err := h.organizationService.UpdateGame(ctx, integratorID, req.GameID, req.WagerSetID, req.RTP, req.Volatility, req.ShortLink)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, ig, nil)
}

// @Summary delete integrator game.
// @Tags organizations
// @Consume application/json
// @Description The provider_id field is required. If you send only the provider_id,
// @Description all games from the provider will be deleted. To delete selected games, fill out the game_id.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Param data body requests.IntegratorGameRequest true "IntegratorGameRequest"
// @Success 204
// @Router /api/organizations/{id}/game [delete].
func (h *organizationHandler) revokeGames(ctx *gin.Context) {
	req := &requests.IntegratorGameRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	integratorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err = h.organizationService.RevokeGames(ctx, integratorID, req.GameID...)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.NoContent(ctx)
}

// @Summary Get organization pair.
// @Tags organizations
// @Consume application/json
// @Description Get organization pair.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.CurrencyOrganizationPair true "CurrencyOrganizationPair"
// @Success 200 {object} response.Response{data=entities.ProviderIntegratorPair}
// @Router /api/organizations/get_organization_pair [post].
func (h *organizationHandler) getOrganizationPair(ctx *gin.Context) {
	req := &requests.CurrencyOrganizationPair{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	pair, err := h.organizationService.GetOrganizationPair(ctx, req.ProviderID, req.IntegratorID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, pair, nil)
}

// @Summary Add organization pair.
// @Tags organizations
// @Consume application/json
// @Description Add organization pair.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.CurrencyOrganizationPair true "CurrencyOrganizationPair"
// @Success 200 {object} response.Response{data=entities.ProviderIntegratorPair}
// @Router /api/organizations/add_organization_pair [post].
func (h *organizationHandler) addOrganizationPair(ctx *gin.Context) {
	req := &requests.CurrencyOrganizationPair{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	pair, err := h.organizationService.CreateOrganizationPair(ctx, req.ProviderID, req.IntegratorID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, pair, nil)
}

// @Summary Delete organization pair.
// @Tags organizations
// @Consume application/json
// @Description Delete organization pair.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.CurrencyOrganizationPair true "CurrencyOrganizationPair"
// @Success 204
// @Router /api/organizations/delete_organization_pair [delete].
func (h *organizationHandler) deleteOrganizationPair(ctx *gin.Context) {
	req := &requests.CurrencyOrganizationPair{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err := h.organizationService.DeleteOrganizationPair(ctx, req.ProviderID, req.IntegratorID)
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

// @Summary Get operator pair.
// @Tags organizations
// @Consume application/json
// @Description Get operator pair.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.IntegratorOperatorPair true "IntegratorOperatorPair"
// @Success 200 {object} response.Response{data=entities.ProviderIntegratorPair}
// @Router /api/organizations/get_operator_pair [post].
func (h *organizationHandler) getOperatorPair(ctx *gin.Context) {
	req := &requests.IntegratorOperatorPair{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	pair, err := h.organizationService.GetOperatorPair(ctx, req.IntegratorID, req.OperatorID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, pair, nil)
}

// @Summary Add operator pair.
// @Tags organizations
// @Consume application/json
// @Description Add operator pair.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.IntegratorOperatorPair true "IntegratorOperatorPair"
// @Success 200 {object} response.Response{data=entities.ProviderIntegratorPair}
// @Router /api/organizations/add_operator_pair [post].
func (h *organizationHandler) addOperatorPair(ctx *gin.Context) {
	req := &requests.IntegratorOperatorPair{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	pair, err := h.organizationService.CreateOperatorPair(ctx, req.IntegratorID, req.OperatorID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, pair, nil)
}

// @Summary Delete operator pair.
// @Tags organizations
// @Consume application/json
// @Description Delete operator pair.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.IntegratorOperatorPair true "IntegratorOperatorPair"
// @Success 204
// @Router /api/organizations/delete_operator_pair [delete].
func (h *organizationHandler) deleteOperatorPair(ctx *gin.Context) {
	req := &requests.IntegratorOperatorPair{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err := h.organizationService.DeleteOperatorPair(ctx, req.IntegratorID, req.OperatorID)
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

// @Summary Get organization operator list by integrator.
// @Tags organizations
// @Consume application/json
// @Description Available operators.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Success 200 {object} response.Response{data=[]entities.Organization}
// @Router /api/organizations/{id}/operators [get].
func (h *organizationHandler) getOperatorsByIntegrator(ctx *gin.Context) {
	integratorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	operators, err := h.organizationService.GetOperatorsByIntegrator(ctx, integratorID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, operators, nil)
}

// @Summary get integrator games wager sets.
// @Tags organizations
// @Consume application/json
// @Description get integrator games wager sets.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Success 200 {object} response.Response{data=[]entities.IntegratorGameWagerSet}
// @Router /api/organizations/{id}/wager_set [get].
func (h *organizationHandler) getGamesWagerSets(ctx *gin.Context) {
	integratorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	igws, err := h.organizationService.GetGamesWagerSets(ctx, integratorID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, igws, nil)
}

// @Summary add integrator game wager sets.
// @Tags organizations
// @Consume application/json
// @Description The provider_id field is required. If you send only the provider_id,
// @Description all games from the provider will be added. To add selected games, fill out the game_id.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Param data body requests.IntegratorGameWagerSetRequest true "IntegratorGameWagerSetRequest"
// @Success 200 {object} response.Response{data=[]entities.IntegratorGameWagerSet}
// @Router /api/organizations/{id}/wager_set [post].
func (h *organizationHandler) addGameWagerSet(ctx *gin.Context) {
	req := &requests.IntegratorGameWagerSetRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	integratorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	zap.S().Info("init method cfgSender.Fire for integrator-games-wager-set")
	igws, err := h.organizationService.CreateIntegratorGameWagerSet(ctx, integratorID, req.WagerSetID, req.Currency, req.GameID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, igws, nil)
}

// @Summary update integrator game wager set.
// @Tags organizations
// @Consume application/json
// @Description The provider_id field is required. If you send only the provider_id,
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Param data body requests.UpdateIntegratorGameWagerSetRequest true "UpdateIntegratorGameWagerSetRequest"
// @Success 200 {object} response.Response{data=entities.IntegratorGameWagerSet}
// @Router /api/organizations/{id}/wager_set [put].
func (h *organizationHandler) updateGameWagerSet(ctx *gin.Context) {
	req := &requests.UpdateIntegratorGameWagerSetRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	integratorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}
	zap.S().Info("init method cfgSender.Fire for integrator-games-wager-set")
	igws, err := h.organizationService.UpdateGameWagerSet(ctx, integratorID, req.GameID, req.WagerSetID, req.Currency, req.NewCurrency)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, igws, nil)
}

// @Summary delete integrator game wager set.
// @Tags organizations
// @Consume application/json
// @Description The provider_id field is required. If you send only the provider_id,
// @Description all games from the provider will be deleted. To delete selected games, fill out the game_id.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "organization_id"
// @Param data body requests.IntegratorGameWagerSetRequest true "IntegratorGameWagerSetRequest"
// @Success 204
// @Router /api/organizations/{id}/wager_set [delete].
func (h *organizationHandler) deleteGameWagerSet(ctx *gin.Context) {
	req := &requests.IntegratorGameWagerSetRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	integratorID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err = h.organizationService.DeleteGameWagerSet(ctx, integratorID, req.WagerSetID, req.Currency, req.GameID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.NoContent(ctx)
}
