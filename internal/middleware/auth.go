package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/go-backend-project/internal/auth"
)

const (
	ContextUserIDKey = "user_id"
	ContextRoleKey   = "role"
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

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil || userID <= 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token is invalid",
			})
			return
		}

		ctx.Set(ContextKeyUserID, userID)
		ctx.Set(ContextKeyRole, claims.Role)

		ctx.Next()
	}
}
