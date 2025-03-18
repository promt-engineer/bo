package middlewares

import (
	"backoffice/internal/services"
	"backoffice/internal/transport/http/response"
	"backoffice/pkg/auth"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func Authenticate(provider auth.Authorizer, sessionService *services.SessionService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jti, err := provider.Verify(ctx)
		if err != nil {
			response.Unauthorized(ctx, err, nil)
			return
		}

		session, err := sessionService.Get(ctx, uuid.MustParse(*jti))
		if err != nil {
			response.Unauthorized(ctx, ErrUnauthorized, nil)
			return
		}

		ctx.Set("session_id", *jti)
		ctx.Set("session", session)
		ctx.Next()
	}
}
