package middlewares

import (
	"backoffice/internal/entities"
	"backoffice/internal/transport/http/response"
	"backoffice/pkg/totp"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	ErrNeedActivateTwoFactor = errors.New("You need activate two factor in settings")
	ErrTOTPSecretRequired    = errors.New("The field TOTP secret is required")
	ErrInvalidTOTPSecret     = errors.New("TOTP secret is invalid")
)

func TOTP(required bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := ctx.Value("session").(*entities.Session)

		if required && !session.Account.TOTPEnabled {
			response.BadRequest(ctx, ErrNeedActivateTwoFactor, nil)
			return
		}

		if session.Account.TOTPEnabled || required {
			req := &totp.Request{}
			if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
				response.Unauthorized(ctx, ErrTOTPSecretRequired, "totp_required")
				return
			}

			valid, err := totp.T().Validate(req.TOTP, session.Account.TOTPSecret)
			if err != nil {
				response.ValidationFailed(ctx, err)
				return
			}

			if !valid {
				response.BadRequest(ctx, ErrInvalidTOTPSecret, nil)
				return
			}
		}

		ctx.Next()
	}
}
