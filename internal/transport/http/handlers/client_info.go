package handlers

import (
	"backoffice/internal/services"
	"backoffice/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

type clientInfoHandler struct {
	clientInfoService *services.ClientInfoService
}

func NewClientInfoHTTPHandler(clientInfoService *services.ClientInfoService) *clientInfoHandler {
	return &clientInfoHandler{
		clientInfoService: clientInfoService,
	}
}

func (h *clientInfoHandler) Register(router *gin.RouterGroup) {
	info := router.Group("info")

	info.GET("environment_name", h.GetEnvironmentName)
	info.GET("logo", h.GetLogoURL)
}

// @Summary Get environment_name.
// @Tags info
// @Consume application/json
// @Description Return backoffice environment name.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Success 200 {object} map[string]interface{}
// @Router /api/info/environment_name [get].
func (h *clientInfoHandler) GetEnvironmentName(ctx *gin.Context) {
	environmentName := h.clientInfoService.GetEnvironmentName()
	response.OK(ctx, environmentName, nil)
}

// @Summary Get logo url
// @Tags info
// @Consume application/json
// @Description Return backoffice logo url.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Success 200 {object} map[string]interface{}
// @Router /api/info/logo [get].
func (h *clientInfoHandler) GetLogoURL(ctx *gin.Context) {
	logoURL := h.clientInfoService.GetLogoURL()
	response.OK(ctx, logoURL, nil)
}
