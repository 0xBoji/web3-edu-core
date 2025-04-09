package repositories

import (
	"github.com/0xBoji/web3-edu-core/internal/database/postgres"
	"github.com/0xBoji/web3-edu-core/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonRepository struct {
	db *gorm.DB
}

// NewLessonRepository creates a new lesson repository
func NewLessonRepository() *LessonRepository {
	return &LessonRepository{
		db: postgres.GetDB(),
	}
}

// Create creates a new lesson
func (r *LessonRepository) Create(lesson *models.Lesson) error {
	return r.db.Create(lesson).Error
}

// GetByID gets a lesson by ID
func (r *LessonRepository) GetByID(id uuid.UUID) (*models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.Where("id = ?", id).First(&lesson).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

// Update updates a lesson
func (r *LessonRepository) Update(lesson *models.Lesson) error {
	return r.db.Save(lesson).Error
}

// Delete deletes a lesson
func (r *LessonRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Lesson{}, id).Error
}

// GetByCourseID gets lessons by course ID
func (r *LessonRepository) GetByCourseID(courseID uuid.UUID) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.Where("course_id = ?", courseID).Order("order_number ASC").Find(&lessons).Error
	if err != nil {
		return nil, err
	}
	return lessons, nil
}

// UpdateOrder updates the order of lessons
func (r *LessonRepository) UpdateOrder(lessons []models.Lesson) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, lesson := range lessons {
			if err := tx.Model(&models.Lesson{}).Where("id = ?", lesson.ID).Update("order_number", lesson.OrderNumber).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
