package handler

import (
	"io"
	"net/http"
	"strconv"

	"lms-system-internship/service"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	service service.AttachmentService
}

func NewAttachmentHandler(s service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: s}
}

func (h *AttachmentHandler) UploadFile(c *gin.Context) {
	lessonIDStr := c.PostForm("lesson_id")
	lessonID, err := strconv.ParseUint(lessonIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson_id"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	attachment, err := h.service.UploadFile(c.Request.Context(), uint(lessonID), header.Filename, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

func (h *AttachmentHandler) DownloadFile(c *gin.Context) {
	idStr := c.Param("attachment_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment_id"})
		return
	}

	data, filename, err := h.service.DownloadFile(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/octet-stream", data)
}
