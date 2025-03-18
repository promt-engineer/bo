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

type accountHandler struct {
	accountService       *services.AccountService
	authorizationService *services.AuthorizationService
	organizationService  *services.OrganizationService
}

func NewAccountHandler(accountService *services.AccountService, authorizationService *services.AuthorizationService, organizationService *services.OrganizationService) *accountHandler {
	return &accountHandler{
		accountService:       accountService,
		authorizationService: authorizationService,
		organizationService:  organizationService,
	}
}

func (h *accountHandler) Register(route *gin.RouterGroup) {
	accounts := route.Group("accounts")
	{
		accounts.GET("", h.all)
		accounts.POST("", h.create)
		account := accounts.Group(":id")
		{
			account.GET("", h.get)
			account.PUT("", h.update)
			account.DELETE("", h.delete)
			account.POST("roles", h.assignRole)
			account.DELETE("roles", h.revokeRole)
			account.POST("change_password", h.changePassword)
			//account.POST("permissions", h.assignPermissions)
			//account.DELETE("permissions", h.revokePermissions)
			account.POST("organizations", h.assignOrganization)
			account.DELETE("organizations", h.revokeOrganization)
			account.POST("operators", h.assignOperator)
			account.DELETE("operators", h.revokeOperator)
		}
	}
}

// @Summary Get account list.
// @Tags accounts
// @Consume application/json
// @Description All accounts.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param limit query int true "rows limit"
// @Param offset query int true "rows offset"
// @Param order query string false "order field"
// @Success 200 {object} response.Response{data=[]entities.Account}
// @Router /api/accounts [get].
func (h *accountHandler) all(ctx *gin.Context) {
	req := &requests.Pagination{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	session := ctx.Value("session").(*entities.Session)
	accounts, total, err := h.accountService.Paginate(ctx, session.OrganizationID, req.Filters, req.Order, req.Limit, req.Offset)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	req.Total = total

	response.OK(ctx, accounts, req)
}

// @Summary Create new account.
// @Tags accounts
// @Consume application/json
// @Description Create account.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.CreateAccountRequest true "CreateAccountRequest"
// @Success 200 {object} response.Response{data=entities.Account}
// @Router /api/accounts [post].
func (h *accountHandler) create(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.CreateAccountRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err := h.authorizationService.CanAssignRole(ctx, req.RoleID)
	if err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	account, err := h.accountService.Create(ctx, req.ID, req.Token, req.FirstName, req.LastName)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	account, err = h.organizationService.Assign(ctx, account.ID, session.OrganizationID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	err = h.authorizationService.AssignRole(ctx, account.ID.String(), req.RoleID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	_, err = h.organizationService.AssignOperator(ctx, account.ID.String(), req.OperatorID, session.OrganizationID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	account, err = h.accountService.FindBy(ctx, map[string]interface{}{"id": account.ID})
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, account, nil)
}

// @Summary Update account.
// @Tags accounts
// @Consume application/json
// @Description update account.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Param data body requests.UpdateAccountRequest true "UpdateAccountRequest"
// @Success 200 {object} response.Response
// @Router /api/accounts/{id} [put].
func (h *accountHandler) update(ctx *gin.Context) {
	req := &requests.UpdateAccountRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	account, err := h.accountService.Update(ctx, ctx.Param("id"), req.FirstName, req.LastName)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)
			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, account, nil)
}

// @Summary Change account password.
// @Tags accounts
// @Consume application/json
// @Description change account password.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Param data body requests.ChangePasswordRequest true "ChangePasswordRequest"
// @Success 200 {object} response.Response
// @Router /api/accounts/{id}/change_password [post].
func (h *accountHandler) changePassword(ctx *gin.Context) {
	req := &requests.ChangePasswordRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	account, err := h.accountService.FindBy(ctx, map[string]interface{}{"id": ctx.Param("id")})
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	err = h.accountService.ChangePassword(ctx, account, req.Password, req.NewPassword)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, "success", nil)
}

// @Summary Delete account.
// @Tags accounts
// @Consume application/json
// @Description delete account.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Success 200 {object} response.Response
// @Router /api/accounts/{id} [delete].
func (h *accountHandler) delete(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	err := h.accountService.Delete(ctx, session.Account, ctx.Param("id"))
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, "Success", nil)
}

// @Summary Add organization.
// @Tags accounts
// @Consume application/json
// @Description Add organization.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Param data body requests.AccountOrganizationRequest true "AccountOrganizationRequest"
// @Success 200 {object} response.Response{data=entities.Account}
// @Router /api/accounts/{id}/organizations [post].
func (h *accountHandler) assignOrganization(ctx *gin.Context) {
	req := &requests.AccountOrganizationRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	account, err := h.organizationService.Assign(ctx, uuid.MustParse(ctx.Param("id")), req.OrganizationID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, account, nil)
}

// @Summary Delete organization.
// @Tags accounts
// @Consume application/json
// @Description Delete organization.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Param data body requests.AccountOrganizationRequest true "AccountOrganizationRequest"
// @Success 200 {object} response.Response
// @Router /api/accounts/{id}/organizations [delete].
func (h *accountHandler) revokeOrganization(ctx *gin.Context) {
	req := &requests.AccountOrganizationRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err := h.organizationService.Revoke(ctx, uuid.MustParse(ctx.Param("id")), req.OrganizationID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, "Success", nil)
}

// @Summary Add role.
// @Tags accounts
// @Consume application/json
// @Description Add role.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Param data body requests.AccountRoleRequest true "AccountRoleRequest"
// @Success 200 {object} response.Response
// @Router /api/accounts/{id}/roles [post].
func (h *accountHandler) assignRole(ctx *gin.Context) {
	req := &requests.AccountRoleRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err := h.authorizationService.AssignRole(ctx, ctx.Param("id"), req.RoleID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, "Success", nil)
}

// @Summary Revoke role.
// @Tags accounts
// @Consume application/json
// @Description Revoke role.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Param data body requests.AccountRoleRequest true "AccountRoleRequest"
// @Success 200 {object} response.Response
// @Router /api/accounts/{id}/roles [delete].
func (h *accountHandler) revokeRole(ctx *gin.Context) {
	req := &requests.AccountRoleRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err := h.authorizationService.RevokeRole(ctx, ctx.Param("id"), req.RoleID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, "Success", nil)
}

// @Summary Add operator.
// @Tags accounts
// @Consume application/json
// @Description Add operator.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Param data body requests.AccountOperatorRequest true "AccountOperatorRequest"
// @Success 200 {object} response.Response{data=entities.Account}
// @Router /api/accounts/{id}/operators [post].
func (h *accountHandler) assignOperator(ctx *gin.Context) {
	req := &requests.AccountOperatorRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	account, err := h.organizationService.AssignOperator(ctx, uuid.MustParse(ctx.Param("id")).String(), req.OperatorID.String(), req.IntegratorID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, account, nil)
}

// @Summary Delete operator.
// @Tags accounts
// @Consume application/json
// @Description Delete operator.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Param data body requests.AccountOperatorRequest true "AccountOperatorRequest"
// @Success 200 {object} response.Response
// @Router /api/accounts/{id}/organizations [delete].
func (h *accountHandler) revokeOperator(ctx *gin.Context) {
	req := &requests.AccountOperatorRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err := h.organizationService.RevokeOperator(ctx, ctx.Param("id"), req.OperatorID, req.IntegratorID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, "Success", nil)
}

//func (h *accountHandler) assignPermissions(ctx *gin.Context) {
//	req := &requests.AccountPermissionsRequest{}
//	if err := ctx.ShouldBind(&req); err != nil {
//		response.ValidationFailed(ctx, err)
//
//		return
//	}
//
//	err := h.authorizationService.AssignAccountPermissions(ctx, ctx.Param("id"), req.Permissions...)
//	if err != nil {
//		response.BadRequest(ctx, err, nil)
//
//		return
//	}
//
//	response.OK(ctx, "Success", nil)
//}

//func (h *accountHandler) revokePermissions(ctx *gin.Context) {
//	req := &requests.AccountPermissionsRequest{}
//	if err := ctx.ShouldBind(&req); err != nil {
//		response.ValidationFailed(ctx, err)
//
//		return
//	}
//
//	err := h.authorizationService.RevokeAccountPermissions(ctx, ctx.Param("id"), req.Permissions...)
//	if err != nil {
//		response.BadRequest(ctx, err, nil)
//
//		return
//	}
//
//	response.OK(ctx, "Success", nil)
//}

// @Summary Get account.
// @Tags accounts
// @Consume application/json
// @Description Get account information.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param id path string true "account_id"
// @Success 200 {object} response.Response{data=entities.Account}
// @Router /api/accounts/{id} [get].
func (h *accountHandler) get(ctx *gin.Context) {
	account, err := h.accountService.FindBy(ctx, map[string]interface{}{"id": ctx.Param("id")})
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)
			return
		}

		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, account, nil)
}
