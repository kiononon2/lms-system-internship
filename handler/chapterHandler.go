package handler

import (
	"fmt"
	"lms-system-internship/entities"
	"lms-system-internship/pkg"
	"lms-system-internship/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChapterHandler struct {
	svc service.ChapterService
}

func NewChapterHandler(svc service.ChapterService) *ChapterHandler {
	return &ChapterHandler{svc: svc}
}

func (h *ChapterHandler) GetAllChapters(c *gin.Context) {
	chapters, err := h.svc.GetAllChapters(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, chapters)
}

func (h *ChapterHandler) GetChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("chapter_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	chapter, err := h.svc.GetChapter(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(pkg.ErrChapterNotFound)
		return
	}
	c.JSON(http.StatusOK, chapter)
}

func (h *ChapterHandler) CreateChapter(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Query("course_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var chapter entities.Chapter
	if err2 := c.ShouldBindJSON(&chapter); err2 != nil {
		fmt.Println(err2)
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err3 := h.svc.AddChapterToCourse(c.Request.Context(), uint(courseID), &chapter); err3 != nil {
		fmt.Println(err3)
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, chapter)
}

func (h *ChapterHandler) UpdateChapterOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("chapter_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var payload struct {
		Order int `json:"order"`
	}
	if err2 := c.ShouldBindJSON(&payload); err2 != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err3 := h.svc.UpdateChapterOrder(c.Request.Context(), uint(id), payload.Order); err3 != nil {
		c.Error(pkg.ErrChapterNotFound)
		return
	}

	c.Status(http.StatusOK)
}

func (h *ChapterHandler) DeleteChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("chapter_id"), 10, 64)
	if err != nil {
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err2 := h.svc.RemoveChapter(c.Request.Context(), uint(id)); err2 != nil {
		c.Error(pkg.ErrChapterNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}
