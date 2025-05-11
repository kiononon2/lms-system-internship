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
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lessons)
}

func (h *LessonHandler) GetLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("lesson_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	lesson, err := h.svc.GetLesson(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(pkg.ErrLessonNotFound)
		return
	}
	c.JSON(http.StatusOK, lesson)
}

func (h *LessonHandler) CreateLesson(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Query("chapter_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var lesson entities.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err := h.svc.AddLessonToChapter(c.Request.Context(), uint(chapterID), &lesson); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, lesson)
}

func (h *LessonHandler) UpdateLessonContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("lesson_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var payload struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err := h.svc.UpdateLessonContent(c.Request.Context(), uint(id), payload.Content); err != nil {
		c.Error(pkg.ErrLessonNotFound)
		return
	}
	c.Status(http.StatusOK)
}

func (h *LessonHandler) ReorderLessons(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("chapter_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err := h.svc.ReorderLessons(c.Request.Context(), uint(chapterID), ids); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("lesson_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err := h.svc.DeleteLesson(c.Request.Context(), uint(id)); err != nil {
		c.Error(pkg.ErrLessonNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}
