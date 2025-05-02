package storage

import (
	"context"
	"errors"
	"lms-system-internship/entities"
)

var (
	ErrNotFound = errors.New("record not found")
)

type CourseRepository interface {
	FindAll(ctx context.Context) ([]*entities.Course, error)
	FindByID(ctx context.Context, id uint) (*entities.Course, error)
	Save(ctx context.Context, course *entities.Course) error
	Update(ctx context.Context, course *entities.Course) error
	Delete(ctx context.Context, id uint) error
}

type ChapterRepository interface {
	FindByCourseID(ctx context.Context, courseID uint) ([]*entities.Chapter, error)
	FindByID(ctx context.Context, id uint) (*entities.Chapter, error)
	Save(ctx context.Context, chapter *entities.Chapter) error
	Update(ctx context.Context, chapter *entities.Chapter) error
	Delete(ctx context.Context, id uint) error
}

type LessonRepository interface {
	FindByChapterID(ctx context.Context, chapterID uint) ([]*entities.Lesson, error)
	FindByID(ctx context.Context, id uint) (*entities.Lesson, error)
	Save(ctx context.Context, lesson *entities.Lesson) error
	Update(ctx context.Context, lesson *entities.Lesson) error
	Delete(ctx context.Context, id uint) error
}

type Repository struct {
	Course  CourseRepository
	Chapter ChapterRepository
	Lesson  LessonRepository
}
