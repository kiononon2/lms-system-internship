package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"lms-system-internship/entities"
	"lms-system-internship/mocks"
	"lms-system-internship/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLessonService_GetAllLessons(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		lessons := []*entities.Lesson{
			{
				ID:          1,
				Name:        "Lesson 1",
				Description: "Description 1",
				Content:     "Content 1",
				Order:       1,
				ChapterID:   1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockRepo.On("FindAll", mock.Anything).Return(lessons, nil)

		service := NewLessonService(mockRepo)
		result, err := service.GetAllLessons(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, lessons, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		mockRepo.On("FindAll", mock.Anything).Return([]*entities.Lesson{}, nil)

		service := NewLessonService(mockRepo)
		result, err := service.GetAllLessons(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		mockRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

		service := NewLessonService(mockRepo)
		result, err := service.GetAllLessons(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestLessonService_GetLesson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		lesson := &entities.Lesson{
			ID:          1,
			Name:        "Lesson 1",
			Description: "Description 1",
			Content:     "Content 1",
			Order:       1,
			ChapterID:   1,
		}

		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(lesson, nil)

		service := NewLessonService(mockRepo)
		result, err := service.GetLesson(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, lesson, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, pkg.ErrLessonNotFound)

		service := NewLessonService(mockRepo)
		result, err := service.GetLesson(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, pkg.ErrLessonNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLessonService_AddLessonToChapter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		lesson := &entities.Lesson{
			Name:        "New Lesson",
			Description: "New Description",
			Content:     "New Content",
			Order:       1,
		}

		mockRepo.On("Save", mock.Anything, lesson).Return(nil)

		service := NewLessonService(mockRepo)
		err := service.AddLessonToChapter(context.Background(), 1, lesson)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), lesson.ChapterID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		lesson := &entities.Lesson{
			Name:        "New Lesson",
			Description: "New Description",
			Content:     "New Content",
			Order:       1,
		}

		mockRepo.On("Save", mock.Anything, lesson).Return(errors.New("database error"))

		service := NewLessonService(mockRepo)
		err := service.AddLessonToChapter(context.Background(), 1, lesson)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLessonService_UpdateLessonContent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		lesson := &entities.Lesson{
			ID:          1,
			Name:        "Lesson 1",
			Description: "Description 1",
			Content:     "Old Content",
			Order:       1,
			ChapterID:   1,
		}

		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(lesson, nil)
		mockRepo.On("Update", mock.Anything, lesson).Return(nil)

		service := NewLessonService(mockRepo)
		err := service.UpdateLessonContent(context.Background(), 1, "New Content")

		assert.NoError(t, err)
		assert.Equal(t, "New Content", lesson.Content)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, pkg.ErrLessonNotFound)

		service := NewLessonService(mockRepo)
		err := service.UpdateLessonContent(context.Background(), 1, "New Content")

		assert.Error(t, err)
		assert.Equal(t, pkg.ErrLessonNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLessonService_ReorderLessons(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		lessons := []*entities.Lesson{
			{
				ID:        1,
				Order:     1,
				ChapterID: 1,
			},
			{
				ID:        2,
				Order:     2,
				ChapterID: 1,
			},
		}

		mockRepo.On("FindByChapterID", mock.Anything, uint(1)).Return(lessons, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Lesson")).Return(nil).Twice()

		service := NewLessonService(mockRepo)
		err := service.ReorderLessons(context.Background(), 1, []uint{2, 1})

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error finding lessons", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		mockRepo.On("FindByChapterID", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

		service := NewLessonService(mockRepo)
		err := service.ReorderLessons(context.Background(), 1, []uint{2, 1})

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error updating lesson", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		lessons := []*entities.Lesson{
			{
				ID:        1,
				Order:     1,
				ChapterID: 1,
			},
		}

		mockRepo.On("FindByChapterID", mock.Anything, uint(1)).Return(lessons, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Lesson")).Return(errors.New("database error"))

		service := NewLessonService(mockRepo)
		err := service.ReorderLessons(context.Background(), 1, []uint{1})

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLessonService_DeleteLesson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

		service := NewLessonService(mockRepo)
		err := service.DeleteLesson(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(pkg.ErrLessonNotFound)

		service := NewLessonService(mockRepo)
		err := service.DeleteLesson(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, pkg.ErrLessonNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.LessonRepository)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("database error"))

		service := NewLessonService(mockRepo)
		err := service.DeleteLesson(context.Background(), 1)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
