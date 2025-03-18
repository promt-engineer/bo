package auth

import (
	"github.com/gin-gonic/gin"
)

type Authorizer interface {
	Verify(ctx *gin.Context) (*string, error)
	Refresh(opts ...TokenOption) (*Auth, error)
	Generate(opts ...GenerateOption) (*Auth, error)
	Token(opts ...TokenOption) (*Token, error)
}
