package handlers

import (
	"backoffice/internal/entities"
	"backoffice/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

type dashboardHandler struct {
}

func NewDashboardHandler() *dashboardHandler {
	return &dashboardHandler{}
}

func (h *dashboardHandler) Register(route *gin.RouterGroup) {
	route.GET("dashboard", h.getIndex)
}

// @Summary Dashboard.
// @Tags dashboard
// @Consume application/json
// @Description Dashboard.
// @Accept  json
// @Produce  json
// @Success 200  {object} entities.Account
// @Router /api/dashboard [get].
func (h *dashboardHandler) getIndex(ctx *gin.Context) {
	account := ctx.Value("account").(*entities.Account)

	response.OK(ctx, account, nil)
}
