package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:100;not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Slug        string    `gorm:"size:100;not null;unique" json:"slug"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
	Courses     []Course  `gorm:"many2many:course_categories;" json:"courses,omitempty"`
}

// TableName specifies the table name for the Category model
func (Category) TableName() string {
	return "categories"
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
