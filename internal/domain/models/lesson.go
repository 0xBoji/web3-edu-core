package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lesson struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CourseID    uuid.UUID `gorm:"type:uuid" json:"course_id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	VideoURL    string    `gorm:"size:255;not null" json:"video_url"`
	VideoID     string    `gorm:"size:100;not null" json:"video_id"`
	Duration    int       `json:"duration,omitempty"` // minutes
	OrderNumber int       `gorm:"not null" json:"order_number"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
}

// TableName specifies the table name for the Lesson model
func (Lesson) TableName() string {
	return "lessons"
}

// BeforeCreate will set a UUID rather than numeric ID
func (l *Lesson) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
