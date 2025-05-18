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

func TestCourseService_GetAllCourses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		courses := []*entities.Course{
			{
				ID:          1,
				Name:        "Course 1",
				Description: "Description 1",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          2,
				Name:        "Course 2",
				Description: "Description 2",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockRepo.On("FindAll", mock.Anything).Return(courses, nil)

		service := NewCourseService(mockRepo)
		result, err := service.GetAllCourses(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, courses, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		mockRepo.On("FindAll", mock.Anything).Return([]*entities.Course{}, nil)

		service := NewCourseService(mockRepo)
		result, err := service.GetAllCourses(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		mockRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

		service := NewCourseService(mockRepo)
		result, err := service.GetAllCourses(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestCourseService_UpdateCourseDetails(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		course := &entities.Course{
			ID:          1,
			Name:        "Updated Course",
			Description: "Updated Description",
		}

		mockRepo.On("Update", mock.Anything, course).Return(nil)

		service := NewCourseService(mockRepo)
		err := service.UpdateCourseDetails(context.Background(), course)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		course := &entities.Course{
			ID:          999,
			Name:        "Non-existent Course",
			Description: "Should fail",
		}

		mockRepo.On("Update", mock.Anything, course).Return(pkg.ErrCourseNotFound)

		service := NewCourseService(mockRepo)
		err := service.UpdateCourseDetails(context.Background(), course)

		assert.Error(t, err)
		assert.Equal(t, pkg.ErrCourseNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		course := &entities.Course{
			ID:          1,
			Name:        "Course",
			Description: "Description",
		}

		mockRepo.On("Update", mock.Anything, course).Return(errors.New("database error"))

		service := NewCourseService(mockRepo)
		err := service.UpdateCourseDetails(context.Background(), course)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCourseService_DeleteCourse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

		service := NewCourseService(mockRepo)
		err := service.DeleteCourse(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(pkg.ErrCourseNotFound)

		service := NewCourseService(mockRepo)
		err := service.DeleteCourse(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, pkg.ErrCourseNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("database error"))

		service := NewCourseService(mockRepo)
		err := service.DeleteCourse(context.Background(), 1)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
