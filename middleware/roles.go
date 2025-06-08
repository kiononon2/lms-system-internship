package middleware

import (
	"lms-system-internship/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireRoles проверяет, содержит ли пользователь хотя бы одну из нужных ролей
func RequireRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("roles")
		if !exists {
			pkg.Logger.Warn("No roles found in context")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No roles in token"})
			return
		}

		userRoles, ok := val.([]string)
		if !ok {
			pkg.Logger.Error("Roles have invalid format in context")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid roles format"})
			return
		}

		pkg.Logger.Debugf("User roles: %v", userRoles)
		pkg.Logger.Debugf("Required roles: %v", requiredRoles)

		for _, required := range requiredRoles {
			for _, userRole := range userRoles {
				if userRole == required {
					pkg.Logger.Infof("Access granted for role: %s", userRole)
					c.Next()
					return
				}
			}
		}

		pkg.Logger.Warnf("Access denied. Required: %v, user has: %v", requiredRoles, userRoles)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: insufficient role " + strings.Join(requiredRoles, ","),
		})
	}
}
