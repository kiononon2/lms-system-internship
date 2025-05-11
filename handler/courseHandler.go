package handler

import (
	"github.com/sirupsen/logrus"
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
		pkg.Logger.WithError(err).Error("Failed to retrieve all courses")
		c.Error(err)
		return
	}
	pkg.Logger.Info("Retrieved all courses")
	c.JSON(http.StatusOK, courses)
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("course_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("course_id", c.Param("course_id")).Error("Invalid course ID format")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	course, err := h.svc.GetCourse(c.Request.Context(), uint(id))
	if err != nil {
		pkg.Logger.WithField("course_id", id).Error("Course not found")
		c.Error(pkg.ErrCourseNotFound)
		return
	}
	pkg.Logger.WithField("course_id", id).Info("Retrieved course details")
	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var course entities.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		pkg.Logger.WithError(err).Error("Failed to bind JSON for new course")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	pkg.Logger.WithFields(logrus.Fields{
		"course_name": course.Name,
		"course_desc": course.Description,
	}).Info("Creating new course")

	if err := h.svc.CreateCourse(c.Request.Context(), &course); err != nil {
		pkg.Logger.WithError(err).Error("Failed to create course")
		c.Error(err)
		return
	}
	pkg.Logger.WithField("course", course).Debug("Course created successfully")
	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("course_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("course_id", c.Param("course_id")).Error("Invalid course ID format")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var course entities.Course
	if err2 := c.ShouldBindJSON(&course); err2 != nil {
		pkg.Logger.WithError(err2).Error("Failed to bind JSON for course update")
		c.Error(pkg.ErrInvalidInput)
		return
	}
	course.ID = uint(id)

	if err3 := h.svc.UpdateCourseDetails(c.Request.Context(), &course); err3 != nil {
		pkg.Logger.WithField("course_id", id).Error("Course update failed")
		c.Error(pkg.ErrCourseNotFound)
		return
	}
	pkg.Logger.WithField("course_id", id).Info("Course updated successfully")
	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("course_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("course_id", c.Param("course_id")).Error("Invalid course ID format")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err2 := h.svc.DeleteCourse(c.Request.Context(), uint(id)); err2 != nil {
		pkg.Logger.WithField("course_id", id).Error("Failed to delete course")
		c.Error(pkg.ErrCourseNotFound)
		return
	}
	pkg.Logger.WithField("course_id", id).Info("Course deleted successfully")
	c.Status(http.StatusNoContent)
}
