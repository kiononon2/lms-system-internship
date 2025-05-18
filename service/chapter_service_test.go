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

func TestChapterService_GetAllChapters(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		chapters := []*entities.Chapter{
			{
				ID:          1,
				Name:        "Chapter 1",
				Description: "Description 1",
				Order:       1,
				CourseID:    1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockRepo.On("FindAll", mock.Anything).Return(chapters, nil)

		service := NewChapterService(mockRepo)
		result, err := service.GetAllChapters(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, chapters, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		mockRepo.On("FindAll", mock.Anything).Return([]*entities.Chapter{}, nil)

		service := NewChapterService(mockRepo)
		result, err := service.GetAllChapters(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		mockRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

		service := NewChapterService(mockRepo)
		result, err := service.GetAllChapters(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestChapterService_GetChapter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		chapter := &entities.Chapter{
			ID:          1,
			Name:        "Chapter 1",
			Description: "Description 1",
			Order:       1,
			CourseID:    1,
		}

		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(chapter, nil)

		service := NewChapterService(mockRepo)
		result, err := service.GetChapter(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, chapter, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, pkg.ErrChapterNotFound)

		service := NewChapterService(mockRepo)
		result, err := service.GetChapter(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, pkg.ErrChapterNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestChapterService_AddChapterToCourse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		chapter := &entities.Chapter{
			Name:        "New Chapter",
			Description: "New Description",
			Order:       1,
		}

		mockRepo.On("Save", mock.Anything, chapter).Return(nil)

		service := NewChapterService(mockRepo)
		err := service.AddChapterToCourse(context.Background(), 1, chapter)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), chapter.CourseID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		chapter := &entities.Chapter{
			Name:        "New Chapter",
			Description: "New Description",
			Order:       1,
		}

		mockRepo.On("Save", mock.Anything, chapter).Return(errors.New("database error"))

		service := NewChapterService(mockRepo)
		err := service.AddChapterToCourse(context.Background(), 1, chapter)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestChapterService_UpdateChapterOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		chapter := &entities.Chapter{
			ID:          1,
			Name:        "Chapter 1",
			Description: "Description 1",
			Order:       1,
			CourseID:    1,
		}

		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(chapter, nil)
		mockRepo.On("Update", mock.Anything, chapter).Return(nil)

		service := NewChapterService(mockRepo)
		err := service.UpdateChapterOrder(context.Background(), 1, 2)

		assert.NoError(t, err)
		assert.Equal(t, 2, chapter.Order)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		mockRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, pkg.ErrChapterNotFound)

		service := NewChapterService(mockRepo)
		err := service.UpdateChapterOrder(context.Background(), 1, 2)

		assert.Error(t, err)
		assert.Equal(t, pkg.ErrChapterNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestChapterService_RemoveChapter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

		service := NewChapterService(mockRepo)
		err := service.RemoveChapter(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(pkg.ErrChapterNotFound)

		service := NewChapterService(mockRepo)
		err := service.RemoveChapter(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, pkg.ErrChapterNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.ChapterRepository)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("database error"))

		service := NewChapterService(mockRepo)
		err := service.RemoveChapter(context.Background(), 1)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
