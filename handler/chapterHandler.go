package handler

import (
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

// GetAllChapters godoc
// @Summary      Get all chapters
// @Description  Retrieves a list of all chapters
// @Tags         chapters
// @Produce      json
// @Success      200  {array}   entities.Chapter
// @Failure      500  {object}  pkg.ErrorResponse
// @Router       /api/chapters [get]
func (h *ChapterHandler) GetAllChapters(c *gin.Context) {
	chapters, err := h.svc.GetAllChapters(c.Request.Context())
	if err != nil {
		pkg.Logger.WithError(err).Error("Failed to retrieve all chapters")
		c.Error(err)
		return
	}
	pkg.Logger.Info("Retrieved all chapters")
	c.JSON(http.StatusOK, chapters)
}

// GetChapter godoc
// @Summary      Get chapter by ID
// @Description  Retrieves a specific chapter by its ID
// @Tags         chapters
// @Produce      json
// @Param        chapter_id  path      int  true  "Chapter ID"
// @Success      200  {object}  entities.Chapter
// @Failure      400  {object}  pkg.ErrorResponse
// @Failure      404  {object}  pkg.ErrorResponse
// @Router       /api/chapters/{chapter_id} [get]
func (h *ChapterHandler) GetChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("chapter_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("chapter_id", c.Param("chapter_id")).Error("Invalid chapter ID format")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	chapter, err := h.svc.GetChapter(c.Request.Context(), uint(id))
	if err != nil {
		pkg.Logger.WithField("chapter_id", id).Error("Chapter not found")
		c.Error(pkg.ErrChapterNotFound)
		return
	}
	pkg.Logger.WithField("chapter_id", id).Info("Retrieved chapter details")
	c.JSON(http.StatusOK, chapter)
}

// CreateChapter godoc
// @Summary      Create a new chapter
// @Description  Adds a new chapter to a specific course
// @Tags         chapters
// @Accept       json
// @Produce      json
// @Param        course_id  query     int               true  "Course ID"
// @Param        chapter    body      entities.Chapter  true  "Chapter data"
// @Success      201  {object}  entities.Chapter
// @Failure      400  {object}  pkg.ErrorResponse
// @Failure      500  {object}  pkg.ErrorResponse
// @Router       /api/chapters [post]
func (h *ChapterHandler) CreateChapter(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Query("course_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("course_id", c.Query("course_id")).Error("Invalid course ID format in query")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var chapter entities.Chapter
	if err2 := c.ShouldBindJSON(&chapter); err2 != nil {
		pkg.Logger.Error("Invalid JSON input while creating chapter")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err3 := h.svc.AddChapterToCourse(c.Request.Context(), uint(courseID), &chapter); err3 != nil {
		pkg.Logger.WithError(err).Error("Failed to add chapter to course")
		c.Error(err)
		return
	}
	pkg.Logger.WithFields(map[string]interface{}{
		"course_id": courseID,
		"chapter":   chapter,
	}).Debug("Chapter created successfully")
	c.JSON(http.StatusCreated, chapter)
}

// UpdateChapterOrder godoc
// @Summary      Update chapter order
// @Description  Updates the order field of a specific chapter
// @Tags         chapters
// @Accept       json
// @Produce      json
// @Param        chapter_id  path      int  true  "Chapter ID"
// @Param        order       body      map[string]int  true  "New order value, e.g. {\"order\": 2}"
// @Success      200
// @Failure      400  {object}  pkg.ErrorResponse
// @Failure      404  {object}  pkg.ErrorResponse
// @Router       /api/chapters/{chapter_id}/order [put]
func (h *ChapterHandler) UpdateChapterOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("chapter_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("chapter_id", c.Param("chapter_id")).Error("Invalid chapter ID format")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	var payload struct {
		Order int `json:"order"`
	}
	if err2 := c.ShouldBindJSON(&payload); err2 != nil {
		pkg.Logger.Error("Invalid JSON input while updating chapter order")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err3 := h.svc.UpdateChapterOrder(c.Request.Context(), uint(id), payload.Order); err3 != nil {
		pkg.Logger.WithField("chapter_id", id).Error("Chapter not found while updating order")
		c.Error(pkg.ErrChapterNotFound)
		return
	}

	pkg.Logger.WithFields(map[string]interface{}{
		"chapter_id": id,
		"new_order":  payload.Order,
	}).Info("Updated chapter order")

	c.Status(http.StatusOK)
}

// DeleteChapter godoc
// @Summary      Delete a chapter
// @Description  Deletes a specific chapter by its ID
// @Tags         chapters
// @Produce      json
// @Param        chapter_id  path      int  true  "Chapter ID"
// @Success      204
// @Failure      400  {object}  pkg.ErrorResponse
// @Failure      404  {object}  pkg.ErrorResponse
// @Router       /api/chapters/{chapter_id} [delete]
func (h *ChapterHandler) DeleteChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("chapter_id"), 10, 64)
	if err != nil {
		pkg.Logger.WithField("chapter_id", c.Param("chapter_id")).Error("Invalid chapter ID format")
		c.Error(pkg.ErrInvalidInput)
		return
	}

	if err2 := h.svc.RemoveChapter(c.Request.Context(), uint(id)); err2 != nil {
		pkg.Logger.WithField("chapter_id", id).Error("Chapter not found while deleting")
		c.Error(pkg.ErrChapterNotFound)
		return
	}
	pkg.Logger.WithField("chapter_id", id).Info("Chapter deleted successfully")
	c.Status(http.StatusNoContent)
}
