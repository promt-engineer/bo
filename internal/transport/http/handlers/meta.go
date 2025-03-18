package handlers

import (
	"backoffice/internal/services"
	"backoffice/internal/transport/http/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var tag = "no tag"

type metaHandler struct {
	debugService *services.DebugService
}

func NewMetaHandler(debugService *services.DebugService) *metaHandler {
	debugService.Subscribe(func(report services.DebugReport) {
		zap.S().Infof("%+v", report)
	})

	return &metaHandler{debugService: debugService}
}

func (h *metaHandler) Register(route *gin.RouterGroup) {
	route.GET("health", h.health)
	route.GET("info", h.info)
	route.GET("debug", h.debug)
}

// @Summary Check health.
// @Tags meta
// @Consume application/json
// @Description Check service health.
// @Accept  json
// @Produce  json
// @Success 200  {object} response.HealthResponse
// @Router /api/health [get].
func (h *metaHandler) health(ctx *gin.Context) {
	response.OK(ctx, response.HealthResponse{Success: "ok"}, nil)
}

// @Summary Check tag.
// @Tags meta
// @Consume application/json
// @Description Check service tag.
// @Accept  json
// @Produce  json
// @Success 200  {object} response.InfoResponse
// @Router /api/info [get].
func (h *metaHandler) info(ctx *gin.Context) {
	response.OK(ctx, response.InfoResponse{Tag: tag}, nil)
}

// @Summary Debug call.
// @Tags meta
// @Consume application/json
// @Description Check debug data.
// @Accept  json
// @Produce  json
// @Success 200  {object} string
// @Router /api/debug [get].
func (h *metaHandler) debug(ctx *gin.Context) {
	err := h.debugService.NotifyAll()
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, "ok", nil)
}
