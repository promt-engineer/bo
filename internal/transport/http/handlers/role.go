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

type roleHandler struct {
	authorizationService  *services.AuthorizationService
	authenticationService *services.AuthenticationService
}

func NewRoleHandler(authorizationService *services.AuthorizationService, authenticationService *services.AuthenticationService) *roleHandler {
	return &roleHandler{
		authorizationService:  authorizationService,
		authenticationService: authenticationService,
	}
}

func (h *roleHandler) Register(route *gin.RouterGroup) {
	roles := route.Group("roles")
	{
		roles.GET("", h.all)
		roles.POST("", h.create)

		role := roles.Group(":id")
		{
			role.GET("", h.get)
			role.PUT("", h.update)
			role.DELETE("", h.delete)
			role.POST("permissions", h.assignPermissions)
			role.DELETE("permissions", h.revokePermissions)
		}
	}
}

// @Summary Get roles list.
// @Tags roles
// @Consume application/json
// @Description Available backoffice roles.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param limit query int true "rows limit"
// @Param offset query int true "rows offset"
// @Param order query string false "order field"
// @Success 200  {object} response.Response{data=[]entities.Role}
// @Router /api/roles [get].
func (h *roleHandler) all(ctx *gin.Context) {
	req := &requests.Pagination{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	session := ctx.Value("session").(*entities.Session)
	roles, total, err := h.authorizationService.PaginateRoles(ctx, session.OrganizationID, req.Filters, req.Order, req.Limit, req.Offset)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	req.Total = total

	response.OK(ctx, roles, req)
}

// @Summary Add new role.
// @Tags roles
// @Consume application/json
// @Description Create role.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.UpsertRoleRequest true "CreateRoleRequest"
// @Success 200 {object} response.Response{data=entities.Role}
// @Router /api/roles [post].
func (h *roleHandler) create(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.UpsertRoleRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	role, err := h.authorizationService.CreateRole(ctx, session, req.Name, req.Description, req.Type)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, role, nil)
}

// @Summary Delete role.
// @Tags roles
// @Consume application/json
// @Description Delete existing role.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "role_id"
// @Success 204
// @Router /api/roles/{id} [delete].
func (h *roleHandler) delete(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	err := h.authorizationService.DeleteRole(ctx, session.OrganizationID, ctx.Param("id"))
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

// @Summary Add permissions.
// @Tags roles
// @Consume application/json
// @Description Add permissions.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "role_id"
// @Param data body requests.RolePermissionsRequest true "RolePermissionsRequest"
// @Success 200 {object} response.Response
// @Router /api/roles/{id}/permissions [post].
func (h *roleHandler) assignPermissions(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.RolePermissionsRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	roleID := ctx.Param("id")

	if err := h.authorizationService.AssignRolePermissions(ctx, session.OrganizationID, roleID, req.Permissions...); err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	if err := h.authenticationService.LogoutAllByRole(ctx, roleID); err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, "Success", nil)
}

// @Summary Revoke permissions.
// @Tags roles
// @Consume application/json
// @Description Delete permissions from role.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "role_id"
// @Param data body requests.RolePermissionsRequest true "RolePermissionsRequest"
// @Success 200 {object} response.Response
// @Router /api/roles/{id}/permissions [delete].
func (h *roleHandler) revokePermissions(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.RolePermissionsRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	roleID := ctx.Param("id")

	if err := h.authorizationService.RevokeRolePermissions(ctx, session, roleID, req.Permissions...); err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	if err := h.authenticationService.LogoutAllByRole(ctx, roleID); err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, "Success", nil)
}

// @Summary Get role.
// @Tags roles
// @Consume application/json
// @Description Get existing role.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "role_id"
// @Success 200 {object} response.Response{data=entities.Role}
// @Router /api/roles/{id} [get].
func (h *roleHandler) get(ctx *gin.Context) {
	roleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	role, err := h.authorizationService.GetRole(ctx, roleID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, role, nil)
}

// @Summary Update role.
// @Tags roles
// @Consume application/json
// @Description Update existing role.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "role_id"
// @Param data body requests.UpsertRoleRequest true "UpdateRoleRequest"
// @Success 200 {object} response.Response{data=entities.Role}
// @Router /api/roles/{id} [put].
func (h *roleHandler) update(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.UpsertRoleRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	roleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	role, err := h.authorizationService.UpdateRole(ctx, session, roleID, req.Name, req.Description, req.Type)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, role, nil)
}
