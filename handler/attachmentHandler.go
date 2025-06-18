package handler

import (
	"github.com/google/uuid"
	"io"
	"lms-system-internship/pkg"
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
	attachmentIDParam := c.Param("attachment_id")
	attachmentID, err := strconv.ParseUint(attachmentIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment ID"})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: userID not found"})
		return
	}
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid userID format"})
		return
	}
	pkg.Logger.Debugf("userID (parsed from context):", userID)
	// Скачиваем файл
	data, fileName, err := h.service.DownloadFile(c.Request.Context(), userID, uint(attachmentID))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Отправляем файл
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.Data(http.StatusOK, "application/octet-stream", data)
}
