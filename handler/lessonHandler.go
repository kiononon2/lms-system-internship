package handler

import (
	"github.com/google/uuid"
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

// GetAllLessons godoc
// @Summary      Get all lessons
// @Description  Retrieves a list of all lessons
// @Tags         lessons
// @Produce      json
// @Success      200  {array}   entities.Lesson
// @Failure      500  {object}  pkg.ErrorResponse
// @Router       /api/lessons [get]
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

// GetLesson godoc
// @Summary      Get lesson by ID
// @Description  Retrieves a specific lesson by its ID
// @Tags         lessons
// @Produce      json
// @Param        lesson_id  path      int  true  "Lesson ID"
// @Success      200  {object}  entities.Lesson
// @Failure      400  {object}  pkg.ErrorResponse
// @Failure      404  {object}  pkg.ErrorResponse
// @Router       /api/lessons/{lesson_id} [get]
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

// CreateLesson godoc
// @Summary      Create a new lesson
// @Description  Adds a new lesson to a specific chapter
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Param        chapter_id  query     int              true  "Chapter ID"
// @Param        lesson      body      entities.Lesson  true  "Lesson data"
// @Success      201  {object}  entities.Lesson
// @Failure      400  {object}  pkg.ErrorResponse
// @Failure      500  {object}  pkg.ErrorResponse
// @Router       /api/lessons [post]
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

// UpdateLessonContent godoc
// @Summary      Update lesson content
// @Description  Updates the content field of a specific lesson
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Param        lesson_id  path  int                      true  "Lesson ID"
// @Param        content    body  map[string]string        true  "New content. Example: {\"content\": \"Updated text\"}"
// @Success      200
// @Failure      400  {object}  pkg.ErrorResponse
// @Failure      404  {object}  pkg.ErrorResponse
// @Router  /api/lessons/{lesson_id} [put]
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

// ReorderLessons godoc
// @Summary      Reorder lessons
// @Description  Reorders the lessons in a chapter based on given list of IDs
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Param        chapter_id  path  int        true  "Chapter ID"
// @Param        ids         body  []uint     true  "New lesson order"
// @Success      200
// @Failure      400  {object}  pkg.ErrorResponse
// @Failure      500  {object}  pkg.ErrorResponse
// @Router       /api/chapters/{chapter_id}/lessons/reorder [put]
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

// DeleteLesson godoc
// @Summary      Delete a lesson
// @Description  Deletes a specific lesson by its ID
// @Tags         lessons
// @Produce      json
// @Param        lesson_id  path  int  true  "Lesson ID"
// @Success      204
// @Failure      400  {object}  pkg.ErrorResponse
// @Failure      404  {object}  pkg.ErrorResponse
// @Router       /api/lessons/{lesson_id} [delete]
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

func (h *LessonHandler) GrantLessonAccess(c *gin.Context) {
	var body struct {
		UserID   string `json:"user_id"`
		LessonID uint   `json:"lesson_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	userUUID, err := uuid.Parse(body.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id UUID"})
		return
	}

	if body.LessonID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lesson_id is required"})
		return
	}

	if err := h.svc.GrantAccess(c.Request.Context(), userUUID, body.LessonID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to grant access"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "access granted"})
}
