package middleware

import (
	"gogroceries/domain"
	"gogroceries/internal/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtAuth helper.JWTInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helper.SendError(c, http.StatusUnauthorized, "Header Authorization needed", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			helper.SendError(c, http.StatusUnauthorized, "Wrong Format header Authorization (Must: Bearer <token>)", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtAuth.ValidateToken(tokenString)
		if err != nil {
			helper.SendError(c, http.StatusUnauthorized, "Token is not valid or expired", err.Error())
			c.Abort()
			return
		}

		c.Set("user_claims", claims)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, exists := c.Get("user_claims")
		if !exists {
			helper.SendError(c, http.StatusForbidden, "User claims not found in context", nil)
			c.Abort()
			return
		}

		claims, ok := userClaims.(*domain.JWTClaims)
		if !ok || !claims.IsAdmin {
			helper.SendError(c, http.StatusForbidden, "Access denied: Only admin", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}