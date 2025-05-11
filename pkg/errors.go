package pkg

import "errors"

var (
	ErrCourseNotFound  = errors.New("course not found")
	ErrChapterNotFound = errors.New("chapter not found")
	ErrLessonNotFound  = errors.New("lesson not found")
	ErrInvalidInput    = errors.New("invalid input")
)
