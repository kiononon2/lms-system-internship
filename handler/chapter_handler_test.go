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

func TestChapterHandler_GetAllChapters(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.ChapterService)
		chapters := []*entities.Chapter{
			{
				ID:          1,
				Name:        "Chapter 1",
				Description: "Description 1",
				CourseID:    1,
			},
			{
				ID:          2,
				Name:        "Chapter 2",
				Description: "Description 2",
				CourseID:    1,
			},
		}

		mockService.On("GetAllChapters", mock.Anything).Return(chapters, nil)

		handler := NewChapterHandler(mockService)
		router := setupRouter()
		router.GET("/api/chapters", handler.GetAllChapters)

		req, _ := http.NewRequest(http.MethodGet, "/api/chapters", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockService := new(mocks.ChapterService)
		mockService.On("GetAllChapters", mock.Anything).Return([]*entities.Chapter{}, nil)

		handler := NewChapterHandler(mockService)
		router := setupRouter()
		router.GET("/api/chapters", handler.GetAllChapters)

		req, _ := http.NewRequest(http.MethodGet, "/api/chapters", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "[]", resp.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.ChapterService)
		mockService.On("GetAllChapters", mock.Anything).Return(nil, errors.New("service error"))

		handler := NewChapterHandler(mockService)
		router := setupRouter()
		router.GET("/api/chapters", handler.GetAllChapters)

		req, _ := http.NewRequest(http.MethodGet, "/api/chapters", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestChapterHandler_GetChapter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.ChapterService)
		chapter := &entities.Chapter{
			ID:          1,
			Name:        "Test Chapter",
			Description: "Test Description",
			Order:       1,
			CourseID:    1,
		}

		mockService.On("GetChapter", mock.Anything, uint(1)).Return(chapter, nil)

		handler := NewChapterHandler(mockService)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/api/chapters/:chapter_id", handler.GetChapter)

		req, _ := http.NewRequest(http.MethodGet, "/api/chapters/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler := NewChapterHandler(nil)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/api/chapters/:chapter_id", handler.GetChapter)

		req, _ := http.NewRequest(http.MethodGet, "/api/chapters/invalid", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.ChapterService)
		mockService.On("GetChapter", mock.Anything, uint(1)).Return(nil, pkg.ErrChapterNotFound)

		handler := NewChapterHandler(mockService)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/api/chapters/:chapter_id", handler.GetChapter)

		req, _ := http.NewRequest(http.MethodGet, "/api/chapters/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestChapterHandler_CreateChapter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.ChapterService)
		chapter := &entities.Chapter{
			Name:        "New Chapter",
			Description: "New Description",
			Order:       1,
		}

		mockService.On("AddChapterToCourse", mock.Anything, uint(1), chapter).Return(nil)

		handler := NewChapterHandler(mockService)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/chapters", handler.CreateChapter)

		body := `{"name":"New Chapter","description":"New Description","order":1}`
		req, _ := http.NewRequest(http.MethodPost, "/api/chapters?course_id=1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		handler := NewChapterHandler(nil)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/chapters", handler.CreateChapter)

		body := `{"invalid json`
		req, _ := http.NewRequest(http.MethodPost, "/api/chapters?course_id=1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("missing course_id", func(t *testing.T) {
		handler := NewChapterHandler(nil)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/chapters", handler.CreateChapter)

		body := `{"name":"New Chapter"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/chapters", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestChapterHandler_UpdateChapterOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.ChapterService)
		mockService.On("UpdateChapterOrder", mock.Anything, uint(1), 2).Return(nil)

		handler := NewChapterHandler(mockService)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.PUT("/api/chapters/:chapter_id/order", handler.UpdateChapterOrder) // Add parameter to route

		body := `{"order":2}`
		req, _ := http.NewRequest(http.MethodPut, "/api/chapters/1/order", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestChapterHandler_DeleteChapter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.ChapterService)
		mockService.On("RemoveChapter", mock.Anything, uint(1)).Return(nil) // Changed to RemoveChapter

		handler := NewChapterHandler(mockService)
		router := setupRouter()
		router.DELETE("/api/chapters/:chapter_id", handler.DeleteChapter)

		req, _ := http.NewRequest(http.MethodDelete, "/api/chapters/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.ChapterService)
		mockService.On("RemoveChapter", mock.Anything, uint(1)).Return(pkg.ErrChapterNotFound) // Changed to RemoveChapter

		handler := NewChapterHandler(mockService)
		router := setupRouter()
		router.DELETE("/api/chapters/:chapter_id", handler.DeleteChapter)

		req, _ := http.NewRequest(http.MethodDelete, "/api/chapters/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler := NewChapterHandler(nil)
		router := setupRouter()
		router.DELETE("/api/chapters/:chapter_id", handler.DeleteChapter)

		req, _ := http.NewRequest(http.MethodDelete, "/api/chapters/invalid", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
