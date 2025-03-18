package handlers

import (
	"backoffice/internal/services"
	"backoffice/internal/transport/http/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type publicReportHandler struct {
	spinService *services.SpinService
}

func NewPublicReportHandler(spinService *services.SpinService) *publicReportHandler {
	return &publicReportHandler{
		spinService: spinService,
	}
}

func (h *publicReportHandler) Register(router *gin.RouterGroup) {
	reports := router.Group("public_reports")

	sessions := reports.Group("sessions")
	sessions.GET(":id", h.session)

	spins := reports.Group("spins")
	spins.GET(":id", h.spin)
}

// @Summary Get public spin information.
// @Tags reports
// @Consume application/json
// @Description For all users.
// @Accept  json
// @Produce  json
// @Param   id path   string true  "session_id"
// @Param   currency query string false "currency"
// @Param   game query string true "game"
// @Success 200  {object} response.Response{data=entities.GamingSession}
// @Router /api/public_reports/sessions/{id} [get].
func (h *publicReportHandler) session(ctx *gin.Context) {

	gamingSession, err := h.spinService.GameSession(ctx, ctx.Query("game"), uuid.MustParse(ctx.Param("id")), ctx.Query("currency"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, gamingSession.Prettify(), nil)
}

// @Summary Get public spin.
// @Tags spins
// @Consume application/json
// @Description Get spin information for all users.
// @Accept  json
// @Produce  json
// @Param   id path   string true  "spin_id"
// @Param   currency query string false "currency"
// @Success 200  {object} response.Response{data=entities.Spin}
// @Router /api/public_reports/spins/{id} [get].
func (h *publicReportHandler) spin(ctx *gin.Context) {

	spin, err := h.spinService.GetSpin(ctx, ctx.Param("id"), ctx.Query("currency"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, spin.Prettify(), nil)
}
