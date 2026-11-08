package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/go-backend-project/internal/auth"
)

func AuthMiddleware(cfg auth.TokenConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Incorrect format: must be 'Bearer <token>",
			})
			return
		}

		tokenString := parts[1]

		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is empty",
			})
			return
		}

		claims, err := auth.ParseAccessToken(cfg, tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token is invalid: " + err.Error(),
			})
			return
		}

		ctx.Set("user_id", claims.Subject)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
