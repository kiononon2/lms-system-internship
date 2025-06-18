package repo

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lms-system-internship/entities"
	"time"
)

type LessonUserRepository interface {
	GrantAccess(userID uuid.UUID, lessonID uint) error
	HasAccess(userID uuid.UUID, lessonID uint) (bool, error)
}

type lessonUserRepository struct {
	db *gorm.DB
}

func NewLessonUserRepository(db *gorm.DB) LessonUserRepository {
	return &lessonUserRepository{db: db}
}

func (r *lessonUserRepository) GrantAccess(userID uuid.UUID, lessonID uint) error {
	var existing entities.LessonUser
	err := r.db.
		Where("user_id = ? AND lesson_id = ?", userID, lessonID).
		First(&existing).Error

	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	access := entities.LessonUser{
		UserID:    userID,
		LessonID:  lessonID,
		CreatedAt: time.Now(),
	}
	return r.db.Create(&access).Error
}

func (r *lessonUserRepository) HasAccess(userID uuid.UUID, lessonID uint) (bool, error) {
	var lu entities.LessonUser
	err := r.db.Where("user_id = ? AND lesson_id = ?", userID, lessonID).First(&lu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}
