package middleware

import (
	"errors"
	"net/http"

	"lms-system-internship/pkg"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Default error response
			status := http.StatusInternalServerError
			message := "Internal server error"

			// Handle specific error types
			switch {
			case errors.Is(err, pkg.ErrCourseNotFound),
				errors.Is(err, pkg.ErrChapterNotFound),
				errors.Is(err, pkg.ErrLessonNotFound):
				status = http.StatusNotFound
				message = err.Error()

			case errors.Is(err, pkg.ErrInvalidInput):
				status = http.StatusBadRequest
				message = err.Error()

			default:
				// For unexpected errors, keep the internal server error status
				// but log the detailed error for debugging
				pkg.Logger.WithError(err).Error("Unexpected error occurred")
			}

			c.JSON(status, pkg.ErrorResponse{
				Message: message,
			})
		}
	}
}
