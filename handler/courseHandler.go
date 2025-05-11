package handler

import (
	"lms-system-internship/entities"
	"lms-system-internship/pkg"
	"lms-system-internship/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	svc service.CourseService
}

func NewCourseHandler(svc service.CourseService) *CourseHandler {
	return &CourseHandler{svc: svc}
}

func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.svc.GetAllCourses(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, courses)
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("course_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	course, err := h.svc.GetCourse(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(pkg.ErrCourseNotFound)
		return
	}
	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var course entities.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err := h.svc.CreateCourse(c.Request.Context(), &course); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("course_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var course entities.Course
	if err2 := c.ShouldBindJSON(&course); err2 != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}
	course.ID = uint(id)

	if err3 := h.svc.UpdateCourseDetails(c.Request.Context(), &course); err3 != nil {
		c.Error(pkg.ErrCourseNotFound)
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("course_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err2 := h.svc.DeleteCourse(c.Request.Context(), uint(id)); err2 != nil {
		c.Error(pkg.ErrCourseNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}
