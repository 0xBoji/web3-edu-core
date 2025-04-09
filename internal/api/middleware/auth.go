package middleware

import (
	"strings"

	"github.com/0xBoji/web3-edu-core/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware for authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c)
			c.Abort()
			return
		}

		// Check if the Authorization header is in the correct format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedResponse(c)
			c.Abort()
			return
		}

		// Parse the token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.UnauthorizedResponse(c)
			c.Abort()
			return
		}

		// Set the user ID and role in the context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware is a middleware for role-based authorization
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the role from the context
		role, exists := c.Get("role")
		if !exists {
			utils.UnauthorizedResponse(c)
			c.Abort()
			return
		}

		// Check if the role is allowed
		roleStr := role.(string)
		for _, r := range roles {
			if r == roleStr {
				c.Next()
				return
			}
		}

		utils.ForbiddenResponse(c)
		c.Abort()
	}
}
