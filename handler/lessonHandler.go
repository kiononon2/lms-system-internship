package handler

import (
	"lms-system-internship/entities"
	"lms-system-internship/pkg"
	"lms-system-internship/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	svc service.LessonService
}

func NewLessonHandler(svc service.LessonService) *LessonHandler {
	return &LessonHandler{svc: svc}
}

func (h *LessonHandler) GetAllLessons(c *gin.Context) {
	lessons, err := h.svc.GetAllLessons(c.Request.Context())
	if err != nil {
		pkg.Logger.WithError(err).Error("Failed to retrieve all lessons")
		c.Error(err)
		return
	}
	pkg.Logger.Info("Retrieved all lessons")
	c.JSON(http.StatusOK, lessons)
}

func (h *LessonHandler) GetLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("lesson_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("lesson_id", c.Param("lesson_id")).Error("Invalid lesson ID format")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	lesson, err := h.svc.GetLesson(c.Request.Context(), uint(id))
	if err != nil {
		pkg.Logger.WithField("lesson_id", id).Error("Lesson not found")
		c.Error(pkg.ErrLessonNotFound)
		return
	}
	pkg.Logger.WithField("lesson_id", id).Info("Retrieved lesson details")
	c.JSON(http.StatusOK, lesson)
}

func (h *LessonHandler) CreateLesson(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Query("chapter_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("chapter_id", c.Query("chapter_id")).Error("Invalid chapter ID format in query")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var lesson entities.Lesson
	if err2 := c.ShouldBindJSON(&lesson); err2 != nil {
		pkg.Logger.Error("Invalid JSON input while creating lesson")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err3 := h.svc.AddLessonToChapter(c.Request.Context(), uint(chapterID), &lesson); err3 != nil {
		pkg.Logger.WithError(err).Error("Failed to add lesson to chapter")
		c.Error(err)
		return
	}
	pkg.Logger.WithFields(map[string]interface{}{
		"chapter_id": chapterID,
		"lesson":     lesson,
	}).Debug("Lesson created successfully")
	c.JSON(http.StatusCreated, lesson)
}

func (h *LessonHandler) UpdateLessonContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("lesson_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("lesson_id", c.Param("lesson_id")).Error("Invalid lesson ID format")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var payload struct {
		Content string `json:"content"`
	}
	if err2 := c.ShouldBindJSON(&payload); err2 != nil {
		pkg.Logger.Error("Invalid JSON input while updating lesson content")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err3 := h.svc.UpdateLessonContent(c.Request.Context(), uint(id), payload.Content); err3 != nil {
		pkg.Logger.WithField("lesson_id", id).Error("Lesson not found while updating order")
		c.Error(pkg.ErrLessonNotFound)
		return
	}
	c.Status(http.StatusOK)
}

func (h *LessonHandler) ReorderLessons(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("chapter_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("chapter_id", c.Query("chapter_id")).Error("Invalid chapter ID format in query")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var ids []uint
	if err2 := c.ShouldBindJSON(&ids); err2 != nil {
		pkg.Logger.Error("Invalid JSON input while reordering lesson content")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err3 := h.svc.ReorderLessons(c.Request.Context(), uint(chapterID), ids); err3 != nil {
		pkg.Logger.WithField("lesson_id", ids).Error(err3)
		c.Error(err3)
		return
	}
	pkg.Logger.WithFields(map[string]interface{}{
		"chapter_id": chapterID,
		"lesson_ids": ids,
	}).Debug("Lessons reordered successfully")
	c.Status(http.StatusOK)
}

func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("lesson_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("lesson_id", c.Param("lesson_id")).Error("Invalid lesson ID format")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err2 := h.svc.DeleteLesson(c.Request.Context(), uint(id)); err2 != nil {
		pkg.Logger.WithField("lesson_id", id).Error("Lesson not found while deleting")
		c.Error(pkg.ErrLessonNotFound)
		return
	}
	pkg.Logger.WithField("lesson_id", id).Info("Lesson deleted successfully")
	c.Status(http.StatusNoContent)
}
