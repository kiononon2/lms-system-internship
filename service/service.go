package service

import (
	"context"
	"lms-system-internship/entities"
	"lms-system-internship/repo"
)

func NewService(repo *repo.Repository) *Service {
	return &Service{
		CourseService:  NewCourseService(repo.Course),
		ChapterService: NewChapterService(repo.Chapter),
		LessonService:  NewLessonService(repo.Lesson),
	}
}

// Course Service Implementation
type courseService struct {
	repo repo.CourseRepository
}

func NewCourseService(repo repo.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) GetAllCourses(ctx context.Context) ([]*entities.Course, error) {
	return s.repo.FindAll(ctx)
}

func (s *courseService) GetCourse(ctx context.Context, courseID uint) (*entities.Course, error) {
	return s.repo.FindByID(ctx, courseID)
}

func (s *courseService) CreateCourse(ctx context.Context, course *entities.Course) error {
	return s.repo.Save(ctx, course)
}

func (s *courseService) UpdateCourseDetails(ctx context.Context, course *entities.Course) error {
	return s.repo.Update(ctx, course)
}

func (s *courseService) DeleteCourse(ctx context.Context, courseID uint) error {
	return s.repo.Delete(ctx, courseID)
}

// Chapter Service Implementation
type chapterService struct {
	repo repo.ChapterRepository
}

func NewChapterService(repo repo.ChapterRepository) ChapterService {
	return &chapterService{repo: repo}
}

func (s *chapterService) GetAllChapters(ctx context.Context) ([]*entities.Chapter, error) {
	return s.repo.FindAll(ctx)
}

func (s *chapterService) GetChapter(ctx context.Context, chapterID uint) (*entities.Chapter, error) {
	return s.repo.FindByID(ctx, chapterID)
}

func (s *chapterService) AddChapterToCourse(ctx context.Context, courseID uint, chapter *entities.Chapter) error {
	chapter.CourseID = courseID
	return s.repo.Save(ctx, chapter)
}

func (s *chapterService) UpdateChapterOrder(ctx context.Context, chapterID uint, newOrder int) error {
	chapter, err := s.repo.FindByID(ctx, chapterID)
	if err != nil {
		return err
	}
	chapter.Order = newOrder
	return s.repo.Update(ctx, chapter)
}

func (s *chapterService) RemoveChapter(ctx context.Context, chapterID uint) error {
	return s.repo.Delete(ctx, chapterID)
}

// Lesson Service Implementation
type lessonService struct {
	repo repo.LessonRepository
}

func NewLessonService(repo repo.LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

func (s *lessonService) GetAllLessons(ctx context.Context) ([]*entities.Lesson, error) {
	return s.repo.FindAll(ctx)
}

func (s *lessonService) GetLesson(ctx context.Context, lessonID uint) (*entities.Lesson, error) {
	return s.repo.FindByID(ctx, lessonID)
}

func (s *lessonService) AddLessonToChapter(ctx context.Context, chapterID uint, lesson *entities.Lesson) error {
	lesson.ChapterID = chapterID
	return s.repo.Save(ctx, lesson)
}

func (s *lessonService) UpdateLessonContent(ctx context.Context, lessonID uint, content string) error {
	lesson, err := s.repo.FindByID(ctx, lessonID)
	if err != nil {
		return err
	}
	lesson.Content = content
	return s.repo.Update(ctx, lesson)
}

func (s *lessonService) ReorderLessons(ctx context.Context, chapterID uint, orderedLessonIDs []uint) error {
	lessons, err := s.repo.FindByChapterID(ctx, chapterID)
	if err != nil {
		return err
	}

	// map lesson IDs to entities
	lessonMap := make(map[uint]*entities.Lesson)
	for _, lesson := range lessons {
		lessonMap[lesson.ID] = lesson
	}

	for order, id := range orderedLessonIDs {
		if lesson, exists := lessonMap[id]; exists {
			lesson.Order = order + 1
			if err := s.repo.Update(ctx, lesson); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *lessonService) DeleteLesson(ctx context.Context, lessonID uint) error {
	return s.repo.Delete(ctx, lessonID)
}
