package handlers

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/services"
	"backoffice/internal/transport/http/middlewares"
	"backoffice/internal/transport/http/requests"
	"backoffice/internal/transport/http/response"
	"backoffice/pkg/auth"
	"backoffice/pkg/totp"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authHandler struct {
	authProvider        auth.Authorizer
	authenticateService *services.AuthenticationService
	accountService      *services.AccountService
	sessionService      *services.SessionService
}

func NewAuthHandler(authProvider auth.Authorizer, authenticateService *services.AuthenticationService, accountService *services.AccountService, sessionService *services.SessionService) *authHandler {
	return &authHandler{
		authProvider:        authProvider,
		authenticateService: authenticateService,
		accountService:      accountService,
		sessionService:      sessionService,
	}
}

func (h *authHandler) Register(route *gin.RouterGroup) {
	auth := route.Group("auth")
	{
		auth.POST("login", h.login)

		auth.POST("password/reset", h.resetPasswordRequest)
		auth.POST("password/reset/:token", h.resetPassword)

		auth.Use(middlewares.Authenticate(h.authProvider, h.sessionService))
		auth.POST("refresh", h.refresh)
		auth.POST("logout", h.logout)
		auth.GET("session", h.session)
		auth.POST("organization", h.switchOrganization)
		auth.POST("password/change", middlewares.TOTP(false), h.changePassword)

		opt := auth.Group("otp")
		{
			opt.POST("generate", middlewares.TOTP(false), h.generateTOTP)
			opt.POST("enable", middlewares.TOTP(false), h.enableTOTP)
			opt.POST("disable", middlewares.TOTP(true), h.disableTOTP)
		}
	}

	route.Use(middlewares.Authenticate(h.authProvider, h.sessionService), middlewares.Authorize())
}

// @Summary Generate TOTP QR.
// @Tags TOTP
// @Consume application/json
// @Description Authenticate account.
// @Accept  json
// @Produce  json
// @Param   data body   totp.Request true  "totp.Request"
// @Success 200  {object} response.Response{data=string}
// @Router /api/auth/totp/generate [post].
func (h *authHandler) generateTOTP(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	key, err := totp.T().Generate(
		totp.NewGenerateOptions(
			totp.WithUsername(session.Account.Email),
		),
	)

	if err != nil {
		response.ServerError(ctx, err, nil)
		return
	}

	image, err := key.Image(300, 300)
	if err != nil {
		response.ServerError(ctx, err, nil)
		return
	}

	if _, err = h.accountService.UpdateTOTPSecret(ctx, session.Account, key.Secret(), key.URL()); err != nil {
		response.ServerError(ctx, err, nil)
		return
	}

	response.OK(ctx, image, nil)
}

// @Summary Generate TOTP QR.
// @Tags TOTP
// @Consume application/json
// @Description Authenticate account.
// @Accept  json
// @Produce  json
// @Param   data body   totp.Request true  "totp.Request"
// @Success 200  {object} response.Response{data=string}
// @Router /api/auth/totp/generate [post].
func (h *authHandler) enableTOTP(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &totp.Request{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.Unauthorized(ctx, middlewares.ErrTOTPSecretRequired, "totp_required")
		return
	}

	valid, err := totp.T().Validate(req.TOTP, session.Account.TOTPSecret)
	if err != nil {
		response.ValidationFailed(ctx, err)
		return
	}

	if !valid {
		response.BadRequest(ctx, middlewares.ErrInvalidTOTPSecret, nil)
		return
	}

	if _, err := h.accountService.EnableTOTP(ctx, session.Account); err != nil {
		response.ServerError(ctx, err, nil)
		return
	}

	response.OK(ctx, "Success", "totp_enabled")
}

func (h *authHandler) disableTOTP(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	if _, err := h.accountService.DisableTOTP(ctx, session.Account); err != nil {
		response.ServerError(ctx, err, nil)
		return
	}

	response.OK(ctx, "Success", "totp_disabled")
}

// @Summary Authenticate.
// @Tags Auth
// @Consume application/json
// @Description Authenticate account.
// @Accept json
// @Produce json
// @Param request body requests.AuthenticateRequest true "Authenticate"
// @Success 200  {object} response.Response{data=auth.Auth}
// @Router /api/auth/login [post].
func (h *authHandler) login(ctx *gin.Context) {
	req := &requests.AuthenticateRequest{}
	if err := ctx.ShouldBindBodyWith(req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)
		return
	}

	account, err := h.accountService.Auth(ctx, req.ID, req.Token)
	if err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	if account.TOTPEnabled {
		r := &totp.Request{}
		if err := ctx.ShouldBindBodyWith(r, binding.JSON); err != nil {
			response.Unauthorized(ctx, middlewares.ErrTOTPSecretRequired, "totp_required")
			return
		}

		valid, err := totp.T().Validate(r.TOTP, account.TOTPSecret)
		if err != nil {
			response.ValidationFailed(ctx, err)
			return
		}

		if !valid {
			response.BadRequest(ctx, middlewares.ErrInvalidTOTPSecret, nil)
			return
		}
	}

	token, err := h.authenticateService.Authenticate(ctx, account)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, services.ErrNotValidPassword) {
			response.Unauthorized(ctx, err, nil)
			return
		}

		response.ServerError(ctx, err, nil)
		return
	}

	response.OK(ctx, token, nil)
}

// @Summary Refresh tokens.
// @Tags Auth
// @Consume application/json
// @Description Refresh tokens.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Success 200  {object} response.Response{data=auth.Auth}
// @Router /api/auth/refresh [post].
func (h *authHandler) refresh(ctx *gin.Context) {
	token, err := h.authenticateService.Refresh(ctx, ctx.GetHeader("X-Authenticate"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, services.ErrNotValidPassword) {
			response.Unauthorized(ctx, err, nil)
			return
		}

		response.ServerError(ctx, err, nil)
		return
	}

	response.OK(ctx, token, nil)
}

// @Summary Logout.
// @Tags Auth
// @Consume application/json
// @Description Logout.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200  {object} nil
// @Router /api/auth/logout [post].
func (h *authHandler) logout(ctx *gin.Context) {
	sessionID := uuid.MustParse(ctx.Value("session_id").(string))

	if err := h.authenticateService.Logout(ctx, sessionID, ctx.GetHeader("X-Authenticate")); err != nil {
		response.ServerError(ctx, err, nil)
		return
	}

	response.OK(ctx, nil, nil)
}

// @Summary Change password.
// @Tags Auth
// @Consume application/json
// @Description Change password.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   data body   requests.ChangePasswordRequest true  "ChangePasswordRequest"
// @Success 200  {object} response.Response{data=string}
// @Router /api/auth/password/change [post].
func (h *authHandler) changePassword(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	req := &requests.ChangePasswordRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)
		return
	}

	if err := h.accountService.ChangePassword(ctx, session.Account, req.Password, req.NewPassword); err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, "Success", "password_changed")
}

// @Summary Session.
// @Tags Auth
// @Consume application/json
// @Description Get account session.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200  {object} response.Response{data=[]entities.Session}
// @Router /api/auth/session [get].
func (h *authHandler) session(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	response.OK(ctx, session.CleanUP(), nil)
}

// @Summary Switch current organization.
// @Tags Auth
// @Consume application/json
// @Description Switch current organization.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param   data body   requests.AccountOrganizationRequest true  "AccountOrganizationRequest"
// @Success 200  {object} response.Response{data=[]entities.Session}
// @Router /api/auth/organization [post].
func (h *authHandler) switchOrganization(ctx *gin.Context) {
	sessionID := uuid.MustParse(ctx.Value("session_id").(string))
	session := ctx.Value("session").(*entities.Session)
	req := &requests.AccountOrganizationRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)
		return
	}

	if err := h.sessionService.SwitchOrganization(ctx, session, sessionID, req.OrganizationID); err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, session, nil)
}

// @Summary Reset password request.
// @Tags Auth
// @Consume application/json
// @Description Reset password request.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param data body requests.ResetPasswordRequest true "ResetPasswordRequest"
// @Success 200 {object} response.Response{data=string}
// @Router /api/auth/password/reset [post].
func (h *authHandler) resetPasswordRequest(ctx *gin.Context) {
	req := &requests.ResetPasswordRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)
		return
	}

	err := h.accountService.ResetPasswordRequest(ctx, req.Email)
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

// @Summary Reset password.
// @Tags Auth
// @Consume application/json
// @Description Reset password.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param token path string true "token"
// @Param data body requests.ResetPassword true "ResetPassword"
// @Success 200 {object} response.Response{data=string}
// @Router /api/auth/password/reset/{token} [post].
func (h *authHandler) resetPassword(ctx *gin.Context) {
	req := &requests.ResetPassword{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)
		return
	}

	err := h.accountService.ResetPassword(ctx, ctx.Param("token"), req.NewPassword)
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
