package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title        string    `gorm:"size:255;not null" json:"title"`
	Description  string    `gorm:"type:text" json:"description,omitempty"`
	Thumbnail    string    `gorm:"size:255" json:"thumbnail,omitempty"`
	InstructorID uuid.UUID `gorm:"type:uuid" json:"instructor_id"`
	Instructor   User      `gorm:"foreignKey:InstructorID" json:"instructor,omitempty"`
	Price        float64   `gorm:"type:decimal(10,2)" json:"price"`
	Level        string    `gorm:"size:50" json:"level,omitempty"` // beginner, intermediate, advanced
	Duration     int       `json:"duration,omitempty"`             // total minutes
	Category     string    `gorm:"size:100" json:"category,omitempty"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:now()" json:"updated_at"`
	Lessons      []Lesson  `gorm:"foreignKey:CourseID" json:"lessons,omitempty"`
}

// TableName specifies the table name for the Course model
func (Course) TableName() string {
	return "courses"
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
