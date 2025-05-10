package repo

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"lms-system-internship/entities"
	"lms-system-internship/storage"
)

//type GormRepository struct {
//	db *gorm.DB
//}

func NewRepository(db *gorm.DB) *storage.Repository {
	return &storage.Repository{
		Course:  &courseRepository{db: db},
		Chapter: &chapterRepository{db: db},
		Lesson:  &lessonRepository{db: db},
	}
}

// Course Repository
type courseRepository struct {
	db *gorm.DB
}

func (r *courseRepository) FindAll(ctx context.Context) ([]*entities.Course, error) {
	var courses []*entities.Course
	err := r.db.WithContext(ctx).Find(&courses).Error
	return courses, err
}

func (r *courseRepository) FindByID(ctx context.Context, id uint) (*entities.Course, error) {
	var course entities.Course
	err := r.db.WithContext(ctx).First(&course, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, storage.ErrNotFound
	}
	return &course, err
}

func (r *courseRepository) Save(ctx context.Context, course *entities.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *courseRepository) Update(ctx context.Context, course *entities.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

func (r *courseRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entities.Course{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return storage.ErrNotFound
	}
	return nil
}

// Chapter Repository
type chapterRepository struct {
	db *gorm.DB
}

func (r *chapterRepository) FindByCourseID(ctx context.Context, courseID uint) ([]*entities.Chapter, error) {
	var chapters []*entities.Chapter
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Order("\"order\"").Find(&chapters).Error
	return chapters, err
}

func (r *chapterRepository) FindByID(ctx context.Context, id uint) (*entities.Chapter, error) {
	var chapter entities.Chapter
	err := r.db.WithContext(ctx).First(&chapter, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, storage.ErrNotFound
	}
	return &chapter, err
}

func (r *chapterRepository) Save(ctx context.Context, chapter *entities.Chapter) error {
	return r.db.WithContext(ctx).Create(chapter).Error
}

func (r *chapterRepository) Update(ctx context.Context, chapter *entities.Chapter) error {
	return r.db.WithContext(ctx).Save(chapter).Error
}

func (r *chapterRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entities.Chapter{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return storage.ErrNotFound
	}
	return nil
}

// Lesson Repository
type lessonRepository struct {
	db *gorm.DB
}

func (r *lessonRepository) FindByChapterID(ctx context.Context, chapterID uint) ([]*entities.Lesson, error) {
	var lessons []*entities.Lesson
	err := r.db.WithContext(ctx).Where("chapter_id = ?", chapterID).Order("\"order\"").Find(&lessons).Error
	return lessons, err
}

func (r *lessonRepository) FindByID(ctx context.Context, id uint) (*entities.Lesson, error) {
	var lesson entities.Lesson
	err := r.db.WithContext(ctx).First(&lesson, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, storage.ErrNotFound
	}
	return &lesson, err
}

func (r *lessonRepository) Save(ctx context.Context, lesson *entities.Lesson) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

func (r *lessonRepository) Update(ctx context.Context, lesson *entities.Lesson) error {
	return r.db.WithContext(ctx).Save(lesson).Error
}

func (r *lessonRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entities.Lesson{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return storage.ErrNotFound
	}
	return nil
}
