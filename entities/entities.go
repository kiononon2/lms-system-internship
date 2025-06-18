package entities

import "time"

type Course struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Chapters []Chapter `gorm:"foreignKey:CourseID" json:"chapters"`
}

type Chapter struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Order       int       `gorm:"not null" json:"order"`
	CourseID    uint      `gorm:"not null" json:"course_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Lessons []Lesson `gorm:"foreignKey:ChapterID" json:"lessons"`
}

type Lesson struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Content     string    `gorm:"type:text" json:"content"`
	Order       int       `gorm:"not null" json:"order"`
	ChapterID   uint      `gorm:"not null" json:"chapter_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Attachment struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(255);not null"`
	URL       string `gorm:"type:varchar(255);not null"`
	LessonID  uint   `gorm:"not null"`
	CreatedAt time.Time
}
