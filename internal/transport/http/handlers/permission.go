package handlers

import (
	e "backoffice/internal/errors"
	"backoffice/internal/services"
	"backoffice/internal/transport/http/requests"
	"backoffice/internal/transport/http/response"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type permissionHandler struct {
	authorizationService *services.AuthorizationService
}

func NewPermissionHandler(authorizationService *services.AuthorizationService) *permissionHandler {
	return &permissionHandler{
		authorizationService: authorizationService,
	}
}

func (h *permissionHandler) Register(route *gin.RouterGroup) {
	permissions := route.Group("permissions")
	{
		permissions.GET("", h.all)
		permissions.POST("", h.create)
		permission := permissions.Group(":id")
		{
			permission.GET("", h.get)
			permission.PUT("", h.update)
			permission.DELETE("", h.delete)
		}
	}
}

// @Summary Get permissions list.
// @Tags permission
// @Consume application/json
// @Description Dashboard.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param limit query int true "rows limit"
// @Param offset query int true "rows offset"
// @Param order query string false "order field"
// @Success 200 {object} response.Response{data=[]entities.Permission}
// @Router /api/permissions [get].
func (h *permissionHandler) all(ctx *gin.Context) {
	req := &requests.Pagination{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	permissions, total, err := h.authorizationService.PaginatePermissions(ctx, req.Filters, req.Order, req.Limit, req.Offset)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	req.Total = total

	response.OK(ctx, permissions, req)
}

// @Summary Create new permission.
// @Tags permission
// @Consume application/json
// @Description Create permission.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.UpsertPermissionRequest true "CreatePermissionRequest"
// @Success 200 {object} response.Response{data=entities.Permission}
// @Router /api/permissions [post].
func (h *permissionHandler) create(ctx *gin.Context) {
	req := &requests.UpsertPermissionRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	permission, err := h.authorizationService.CreatePermission(ctx, req.Name, req.Description, req.Subject, req.Endpoint, req.Action)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, permission, nil)
}

// @Summary Get permission.
// @Tags permission
// @Consume application/json
// @Description Get permission.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "permission_id"
// @Success 200 {object} response.Response{data=entities.Permission}
// @Router /api/permissions/{id} [get].
func (h *permissionHandler) get(ctx *gin.Context) {
	permissionID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	permission, err := h.authorizationService.GetPermission(ctx, permissionID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)
			return
		}

		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, permission, nil)
}

// @Summary Update new permission.
// @Tags permission
// @Consume application/json
// @Description Update permission.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "permission_id"
// @Param data body requests.UpsertPermissionRequest true "UpdatePermissionRequest"
// @Success 200 {object} response.Response{data=entities.Permission}
// @Router /api/permissions/{id} [put].
func (h *permissionHandler) update(ctx *gin.Context) {
	req := &requests.UpsertPermissionRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	permissionID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	permission, err := h.authorizationService.UpdatePermission(ctx, permissionID, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, permission, nil)
}

// @Summary Delete permission.
// @Tags permission
// @Consume application/json
// @Description Delete permission.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "permission_id"
// @Success 204
// @Router /api/permissions/{id} [delete].
func (h *permissionHandler) delete(ctx *gin.Context) {
	permissionID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	err = h.authorizationService.DeletePermission(ctx, permissionID)
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
