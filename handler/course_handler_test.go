package handler

import (
	"bytes"
	"errors"
	"lms-system-internship/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"lms-system-internship/entities"
	"lms-system-internship/mocks"
	"lms-system-internship/pkg"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler()) // Add this line
	return r
}

func TestCourseHandler_GetAllCourses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		courses := []*entities.Course{
			{
				ID:          1,
				Name:        "Course 1",
				Description: "Description 1",
			},
			{
				ID:          2,
				Name:        "Course 2",
				Description: "Description 2",
			},
		}

		mockService.On("GetAllCourses", mock.Anything).Return(courses, nil)

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.GET("/api/courses", handler.GetAllCourses)

		req, _ := http.NewRequest(http.MethodGet, "/api/courses", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		mockService.On("GetAllCourses", mock.Anything).Return(nil, errors.New("service error"))

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.GET("/api/courses", handler.GetAllCourses)

		req, _ := http.NewRequest(http.MethodGet, "/api/courses", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCourseHandler_GetCourse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		course := &entities.Course{
			ID:          1,
			Name:        "Test Course",
			Description: "Test Description",
		}

		mockService.On("GetCourse", mock.Anything, uint(1)).Return(course, nil)

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.GET("/api/courses/:course_id", handler.GetCourse)

		req, _ := http.NewRequest(http.MethodGet, "/api/courses/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler := NewCourseHandler(nil)
		router := setupRouter()
		router.GET("/api/courses/:course_id", handler.GetCourse)

		req, _ := http.NewRequest(http.MethodGet, "/api/courses/invalid", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		mockService.On("GetCourse", mock.Anything, uint(1)).Return(nil, pkg.ErrCourseNotFound)

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.GET("/api/courses/:course_id", handler.GetCourse)

		req, _ := http.NewRequest(http.MethodGet, "/api/courses/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCourseHandler_CreateCourse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		course := &entities.Course{
			Name:        "New Course",
			Description: "New Description",
		}

		mockService.On("CreateCourse", mock.Anything, course).Return(nil)

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.POST("/api/courses", handler.CreateCourse)

		body := `{"name":"New Course","description":"New Description"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/courses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		handler := NewCourseHandler(nil)
		router := setupRouter()
		router.POST("/api/courses", handler.CreateCourse)

		body := `{"invalid json`
		req, _ := http.NewRequest(http.MethodPost, "/api/courses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		course := &entities.Course{
			Name:        "New Course",
			Description: "New Description",
		}

		mockService.On("CreateCourse", mock.Anything, course).Return(errors.New("service error"))

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.POST("/api/courses", handler.CreateCourse)

		body := `{"name":"New Course","description":"New Description"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/courses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCourseHandler_UpdateCourse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		course := &entities.Course{
			ID:          1,
			Name:        "Updated Course",
			Description: "Updated Description",
		}

		mockService.On("UpdateCourseDetails", mock.Anything, course).Return(nil)

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.PUT("/api/courses/:course_id", handler.UpdateCourse)

		body := `{"id":1,"name":"Updated Course","description":"Updated Description"}`
		req, _ := http.NewRequest(http.MethodPut, "/api/courses/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler := NewCourseHandler(nil)
		router := setupRouter()
		router.PUT("/api/courses/:course_id", handler.UpdateCourse)

		req, _ := http.NewRequest(http.MethodPut, "/api/courses/invalid", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		course := &entities.Course{
			ID:          1,
			Name:        "Non-existent Course",
			Description: "Should fail",
		}

		mockService.On("UpdateCourseDetails", mock.Anything, course).Return(pkg.ErrCourseNotFound)

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.PUT("/api/courses/:course_id", handler.UpdateCourse)

		body := `{"id":1,"name":"Non-existent Course","description":"Should fail"}`
		req, _ := http.NewRequest(http.MethodPut, "/api/courses/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		handler := NewCourseHandler(nil)
		router := setupRouter()
		router.PUT("/api/courses/:course_id", handler.UpdateCourse)

		body := `{"invalid json`
		req, _ := http.NewRequest(http.MethodPut, "/api/courses/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("id mismatch", func(t *testing.T) {
		handler := NewCourseHandler(nil)
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.PUT("/api/courses/:course_id", handler.UpdateCourse)

		body := `{"id":2,"name":"Course","description":"Mismatched ID"}`
		req, _ := http.NewRequest(http.MethodPut, "/api/courses/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "invalid input")
	})
}

func TestCourseHandler_DeleteCourse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		mockService.On("DeleteCourse", mock.Anything, uint(1)).Return(nil)

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.DELETE("/api/courses/:course_id", handler.DeleteCourse)

		req, _ := http.NewRequest(http.MethodDelete, "/api/courses/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.CourseService)
		mockService.On("DeleteCourse", mock.Anything, uint(1)).Return(pkg.ErrCourseNotFound)

		handler := NewCourseHandler(mockService)
		router := setupRouter()
		router.DELETE("/api/courses/:course_id", handler.DeleteCourse)

		req, _ := http.NewRequest(http.MethodDelete, "/api/courses/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler := NewCourseHandler(nil)
		router := setupRouter()
		router.DELETE("/api/courses/:course_id", handler.DeleteCourse)

		req, _ := http.NewRequest(http.MethodDelete, "/api/courses/invalid", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
