package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRoles проверяет, содержит ли пользователь хотя бы одну из нужных ролей
func RequireRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No roles in token"})
			return
		}

		userRoles, ok := val.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid roles format"})
			return
		}

		for _, required := range requiredRoles {
			for _, userRole := range userRoles {
				if userRole == required {
					print(userRole) //---------
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient role"})
	}
}
