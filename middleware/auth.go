package middleware

import (
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// TokenAuthMiddleware проверяет JWT токен и сохраняет роли пользователя в контексте
func TokenAuthMiddleware(jwks *keyfunc.JWKS) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Извлекаем имя пользователя
		username, _ := claims["preferred_username"].(string)
		c.Set("username", username)

		// Извлекаем роли из resource_access
		roles := []string{}
		if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
			if clientAccess, ok := resourceAccess["lms-app"].(map[string]interface{}); ok {
				if rolesRaw, ok := clientAccess["roles"].([]interface{}); ok {
					for _, role := range rolesRaw {
						if roleStr, ok := role.(string); ok {
							roles = append(roles, roleStr)
						}
					}
				}
			}
		}
		c.Set("roles", roles)

		c.Next()
	}
}
