package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"lms-system-internship/entities"
	"lms-system-internship/middleware"
	"lms-system-internship/mocks"
	"lms-system-internship/pkg"
)

func TestLessonHandler_GetAllLessons(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		lessons := []*entities.Lesson{
			{
				ID:          1,
				Name:        "Lesson 1",
				Description: "Description 1",
				ChapterID:   1,
			},
			{
				ID:          2,
				Name:        "Lesson 2",
				Description: "Description 2",
				ChapterID:   1,
			},
		}

		mockService.On("GetAllLessons", mock.Anything).Return(lessons, nil)

		handler := NewLessonHandler(mockService)
		router := setupRouter()
		router.GET("/api/lessons", handler.GetAllLessons)

		req, _ := http.NewRequest(http.MethodGet, "/api/lessons", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		mockService.On("GetAllLessons", mock.Anything).Return([]*entities.Lesson{}, nil)

		handler := NewLessonHandler(mockService)
		router := setupRouter()
		router.GET("/api/lessons", handler.GetAllLessons)

		req, _ := http.NewRequest(http.MethodGet, "/api/lessons", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "[]", resp.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		mockService.On("GetAllLessons", mock.Anything).Return(nil, errors.New("service error"))

		handler := NewLessonHandler(mockService)
		router := setupRouter()
		router.GET("/api/lessons", handler.GetAllLessons)

		req, _ := http.NewRequest(http.MethodGet, "/api/lessons", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestLessonHandler_GetLesson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		lesson := &entities.Lesson{
			ID:          1,
			Name:        "Test Lesson",
			Description: "Test Description",
			Content:     "Test Content",
			Order:       1,
			ChapterID:   1,
		}

		mockService.On("GetLesson", mock.Anything, uint(1)).Return(lesson, nil)

		handler := NewLessonHandler(mockService)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/api/lessons/:lesson_id", handler.GetLesson)

		req, _ := http.NewRequest(http.MethodGet, "/api/lessons/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler := NewLessonHandler(nil)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/api/lessons/:lesson_id", handler.GetLesson)

		req, _ := http.NewRequest(http.MethodGet, "/api/lessons/invalid", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		mockService.On("GetLesson", mock.Anything, uint(1)).Return(nil, pkg.ErrLessonNotFound)

		handler := NewLessonHandler(mockService)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/api/lessons/:lesson_id", handler.GetLesson)

		req, _ := http.NewRequest(http.MethodGet, "/api/lessons/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestLessonHandler_CreateLesson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		lesson := &entities.Lesson{
			Name:        "New Lesson",
			Description: "New Description",
			Content:     "New Content",
			Order:       1,
		}

		mockService.On("AddLessonToChapter", mock.Anything, uint(1), lesson).Return(nil)

		handler := NewLessonHandler(mockService)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/lessons", handler.CreateLesson)

		body := `{"name":"New Lesson","description":"New Description","content":"New Content","order":1}`
		req, _ := http.NewRequest(http.MethodPost, "/api/lessons?chapter_id=1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		handler := NewLessonHandler(nil)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/lessons", handler.CreateLesson)

		body := `{"invalid json`
		req, _ := http.NewRequest(http.MethodPost, "/api/lessons?chapter_id=1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("missing chapter_id", func(t *testing.T) {
		handler := NewLessonHandler(nil)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/lessons", handler.CreateLesson)

		body := `{"name":"New Lesson"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/lessons", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestLessonHandler_UpdateLessonContent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		mockService.On("UpdateLessonContent", mock.Anything, uint(1), "New Content").Return(nil)

		handler := NewLessonHandler(mockService)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.PUT("/api/lessons/:lesson_id", handler.UpdateLessonContent) // Add parameter to route

		body := `{"content":"New Content"}`
		req, _ := http.NewRequest(http.MethodPut, "/api/lessons/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestLessonHandler_ReorderLessons(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		lessonIDs := []uint{2, 1, 3}
		mockService.On("ReorderLessons", mock.Anything, uint(1), lessonIDs).Return(nil)

		handler := NewLessonHandler(mockService)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.PUT("/api/chapters/:chapter_id/lessons/reorder", handler.ReorderLessons) // Add parameter to route

		body := `[2,1,3]`
		req, _ := http.NewRequest(http.MethodPut, "/api/chapters/1/lessons/reorder", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestLessonHandler_DeleteLesson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		mockService.On("DeleteLesson", mock.Anything, uint(1)).Return(nil)

		handler := NewLessonHandler(mockService)
		router := setupRouter()
		router.DELETE("/api/lessons/:lesson_id", handler.DeleteLesson)

		req, _ := http.NewRequest(http.MethodDelete, "/api/lessons/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.LessonService)
		mockService.On("DeleteLesson", mock.Anything, uint(1)).Return(pkg.ErrLessonNotFound)

		handler := NewLessonHandler(mockService)
		router := setupRouter()
		router.DELETE("/api/lessons/:lesson_id", handler.DeleteLesson)

		req, _ := http.NewRequest(http.MethodDelete, "/api/lessons/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler := NewLessonHandler(nil)
		router := setupRouter()
		router.DELETE("/api/lessons/:lesson_id", handler.DeleteLesson)

		req, _ := http.NewRequest(http.MethodDelete, "/api/lessons/invalid", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
