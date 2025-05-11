package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if any errors occurred during the request
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// You can customize error handling here based on error types or content
			status := http.StatusInternalServerError
			message := err.Error()

			// Example: Check if it's a "not found" error message
			if message == "course not found" || message == "chapter not found" || message == "lesson not found" {
				status = http.StatusNotFound
			}

			c.JSON(status, gin.H{
				"error": message,
			})
		}
	}
}
