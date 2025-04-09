package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Progress struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LessonID       uuid.UUID `gorm:"type:uuid" json:"lesson_id"`
	Lesson         Lesson    `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	PositionSeconds int       `gorm:"default:0" json:"position_seconds"`
	Completed      bool      `gorm:"default:false" json:"completed"`
	LastWatchedAt  time.Time `gorm:"default:now()" json:"last_watched_at"`
}

// TableName specifies the table name for the Progress model
func (Progress) TableName() string {
	return "progress"
}

// BeforeCreate will set a UUID rather than numeric ID
func (p *Progress) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
