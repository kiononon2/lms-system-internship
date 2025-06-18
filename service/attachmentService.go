package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"lms-system-internship/entities"
	"lms-system-internship/files"
	"lms-system-internship/repo"
	"path/filepath"
)

type AttachmentService interface {
	UploadFile(ctx context.Context, lessonID uint, fileName string, fileBytes []byte) (*entities.Attachment, error)
	DownloadFile(ctx context.Context, attachmentID uint) ([]byte, string, error)
	GetAttachmentsByLesson(ctx context.Context, lessonID uint) ([]*entities.Attachment, error)
}

type attachmentService struct {
	repo        repo.AttachmentRepository
	lessonRepo  repo.LessonRepository
	fileStorage files.FileStorage
}

func NewAttachmentService(
	repo repo.AttachmentRepository,
	lessonRepo repo.LessonRepository,
	fileStorage files.FileStorage,
) AttachmentService {
	return &attachmentService{
		repo:        repo,
		lessonRepo:  lessonRepo,
		fileStorage: fileStorage,
	}
}

func (s *attachmentService) UploadFile(ctx context.Context, lessonID uint, fileName string, fileBytes []byte) (*entities.Attachment, error) {
	_, err := s.lessonRepo.FindByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(fileName)         // Например, ".pdf"
	safeName := uuid.New().String() + ext // UUID.ext
	_, err = s.fileStorage.UploadFile(ctx, safeName, fileBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	attachment := &entities.Attachment{
		Name:     fileName,
		URL:      safeName, // Сохраняем только UUID.ext
		LessonID: lessonID,
	}
	if err := s.repo.Save(ctx, attachment); err != nil {
		return nil, fmt.Errorf("failed to save attachment: %w", err)
	}

	return attachment, nil
}

func (s *attachmentService) DownloadFile(ctx context.Context, attachmentID uint) ([]byte, string, error) {
	attachment, err := s.repo.FindByID(ctx, attachmentID)
	if err != nil {
		return nil, "", err
	}

	data, err := s.fileStorage.DownloadFile(ctx, attachment.URL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download from storage: %w", err)
	}

	return data, attachment.Name, nil
}
func (s *attachmentService) GetAttachmentsByLesson(ctx context.Context, lessonID uint) ([]*entities.Attachment, error) {
	return s.repo.FindByLessonID(ctx, lessonID)
}
