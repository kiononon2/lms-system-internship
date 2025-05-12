package pkg

import "errors"

type ErrorResponse struct {
	Message string `json:"error"` // JSON ключ — "error", как у тебя и было
}

var (
	ErrCourseNotFound  = errors.New("course not found")
	ErrChapterNotFound = errors.New("chapter not found")
	ErrLessonNotFound  = errors.New("lesson not found")
	ErrInvalidInput    = errors.New("invalid input")
)
