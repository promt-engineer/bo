package middlewares

import (
	"backoffice/internal/entities"
	"backoffice/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := ctx.Value("session").(*entities.Session)
		if !session.Account.Authorized(ctx.FullPath(), ctx.Request.Method) {
			response.Unauthorized(ctx, "You don't have permission for this action", nil)
			return
		}

		ctx.Next()
	}
}
