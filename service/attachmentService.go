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
	DownloadFile(ctx context.Context, userID uuid.UUID, attachmentID uint) ([]byte, string, error)
	GetAttachmentsByLesson(ctx context.Context, lessonID uint) ([]*entities.Attachment, error)
	//GrantAccess(ctx context.Context, userID uuid.UUID, lessonID uint) error
}

type attachmentService struct {
	repo           repo.AttachmentRepository
	lessonRepo     repo.LessonRepository
	lessonUserRepo repo.LessonUserRepository
	fileStorage    files.FileStorage
}

func NewAttachmentService(repo repo.AttachmentRepository, lessonRepo repo.LessonRepository, lessonUserRepo repo.LessonUserRepository, fileStorage files.FileStorage) *attachmentService {
	return &attachmentService{
		repo:           repo,
		lessonRepo:     lessonRepo,
		lessonUserRepo: lessonUserRepo,
		fileStorage:    fileStorage,
	}
}

func (s *attachmentService) UploadFile(ctx context.Context, lessonID uint, fileName string, fileBytes []byte) (*entities.Attachment, error) {
	_, err := s.lessonRepo.FindByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(fileName)
	safeName := uuid.New().String() + ext

	_, err = s.fileStorage.UploadFile(ctx, safeName, fileBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	attachment := &entities.Attachment{
		Name:     fileName,
		URL:      safeName,
		LessonID: lessonID,
	}
	if err := s.repo.Save(ctx, attachment); err != nil {
		return nil, fmt.Errorf("failed to save attachment: %w", err)
	}

	return attachment, nil
}

func (s *attachmentService) DownloadFile(ctx context.Context, userID uuid.UUID, attachmentID uint) ([]byte, string, error) {
	// Получаем attachment
	attachment, err := s.repo.FindByID(ctx, attachmentID)
	if err != nil {
		return nil, "", err
	}

	// Проверка доступа пользователя к уроку
	hasAccess, err := s.lessonUserRepo.HasAccess(userID, attachment.LessonID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check lesson access: %w", err)
	}
	if !hasAccess {
		return nil, "", fmt.Errorf("access denied: no permission to download this lesson's file")
	}

	// Скачиваем файл
	data, err := s.fileStorage.DownloadFile(ctx, attachment.URL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download from storage: %w", err)
	}

	return data, attachment.Name, nil
}

func (s *attachmentService) GetAttachmentsByLesson(ctx context.Context, lessonID uint) ([]*entities.Attachment, error) {
	return s.repo.FindByLessonID(ctx, lessonID)
}

//func (s *attachmentService) GrantAccess(ctx context.Context, userID uuid.UUID, lessonID uint) error {
//	return s.lessonUserRepo.GrantAccess(userID, lessonID)
//}
