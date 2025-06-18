package repo

import (
	"context"
	"lms-system-internship/entities"

	"gorm.io/gorm"
)

type AttachmentRepository interface {
	Save(ctx context.Context, a *entities.Attachment) error
	FindByID(ctx context.Context, id uint) (*entities.Attachment, error)
	FindByLessonID(ctx context.Context, lessonID uint) ([]*entities.Attachment, error)
}

type attachmentRepo struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepo{db}
}

func (r *attachmentRepo) Save(ctx context.Context, a *entities.Attachment) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *attachmentRepo) FindByID(ctx context.Context, id uint) (*entities.Attachment, error) {
	var a entities.Attachment
	err := r.db.WithContext(ctx).First(&a, id).Error
	return &a, err
}

func (r *attachmentRepo) FindByLessonID(ctx context.Context, lessonID uint) ([]*entities.Attachment, error) {
	var list []*entities.Attachment
	err := r.db.WithContext(ctx).Where("lesson_id = ?", lessonID).Find(&list).Error
	return list, err
}
