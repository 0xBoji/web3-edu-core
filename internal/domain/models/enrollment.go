package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Enrollment struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CourseID   uuid.UUID `gorm:"type:uuid" json:"course_id"`
	Course     Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	EnrolledAt time.Time `gorm:"default:now()" json:"enrolled_at"`
}

// TableName specifies the table name for the Enrollment model
func (Enrollment) TableName() string {
	return "enrollments"
}

// BeforeCreate will set a UUID rather than numeric ID
func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
