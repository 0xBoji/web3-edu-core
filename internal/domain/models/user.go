package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email          string    `gorm:"size:255;not null;unique" json:"email"`
	PasswordHash   string    `gorm:"size:255;not null" json:"-"`
	FullName       string    `gorm:"size:255;not null" json:"full_name"`
	Role           string    `gorm:"size:50;not null;default:user" json:"role"`
	ProfilePicture string    `gorm:"size:255" json:"profile_picture,omitempty"`
	CreatedAt      time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:now()" json:"updated_at"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
