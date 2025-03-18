package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*", "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{"*", "Accept", "Authorization", "Content-Type",
			"X-CSRF-Token", "x-auth", "Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods", "Access-Control-Allow-Credentials"},
		ExposeHeaders:    []string{"*", "Link"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
