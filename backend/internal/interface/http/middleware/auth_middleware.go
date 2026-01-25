package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
)

// AuthMiddleware creates middleware that validates JWT tokens
func AuthMiddleware(authUsecase usecase.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check for Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token and get user info
		userInfo, err := authUsecase.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user", *userInfo)
		c.Set("user_id", userInfo.ID)
		c.Set("user_role", userInfo.Role)

		c.Next()
	}
}

// RequireRole creates middleware that checks if user has required role
func RequireRole(allowedRoles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user := userInfo.(domain.UserInfo)

		// Check if user's role is in allowed roles
		hasRole := false
		for _, role := range allowedRoles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin is a shorthand for RequireRole(domain.RoleAdmin)
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(domain.RoleAdmin)
}

// RequireManagerOrAdmin is a shorthand for RequireRole(domain.RoleAdmin, domain.RoleManager)
func RequireManagerOrAdmin() gin.HandlerFunc {
	return RequireRole(domain.RoleAdmin, domain.RoleManager)
}
