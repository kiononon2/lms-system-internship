package service

import (
	"context"
	"github.com/google/uuid"
	"lms-system-internship/entities"
)

type CourseService interface {
	GetAllCourses(ctx context.Context) ([]*entities.Course, error)
	GetCourse(ctx context.Context, courseID uint) (*entities.Course, error)
	CreateCourse(ctx context.Context, course *entities.Course) error
	UpdateCourseDetails(ctx context.Context, course *entities.Course) error
	DeleteCourse(ctx context.Context, courseID uint) error
}

type ChapterService interface {
	GetAllChapters(ctx context.Context) ([]*entities.Chapter, error)
	GetChapter(ctx context.Context, chapterID uint) (*entities.Chapter, error)
	AddChapterToCourse(ctx context.Context, courseID uint, chapter *entities.Chapter) error
	UpdateChapterOrder(ctx context.Context, chapterID uint, newOrder int) error
	RemoveChapter(ctx context.Context, chapterID uint) error
}

type LessonService interface {
	GetAllLessons(ctx context.Context) ([]*entities.Lesson, error)
	GetLesson(ctx context.Context, lessonID uint) (*entities.Lesson, error)
	AddLessonToChapter(ctx context.Context, chapterID uint, lesson *entities.Lesson) error
	UpdateLessonContent(ctx context.Context, lessonID uint, content string) error
	ReorderLessons(ctx context.Context, chapterID uint, orderedLessonIDs []uint) error
	DeleteLesson(ctx context.Context, lessonID uint) error
	GrantAccess(ctx context.Context, userID uuid.UUID, lessonID uint) error
}

type Service struct {
	CourseService     CourseService
	ChapterService    ChapterService
	LessonService     LessonService
	AttachmentService AttachmentService
}
