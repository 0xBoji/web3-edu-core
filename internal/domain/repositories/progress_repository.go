package repositories

import (
	"time"

	"github.com/0xBoji/web3-edu-core/internal/database/postgres"
	"github.com/0xBoji/web3-edu-core/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProgressRepository struct {
	db *gorm.DB
}

// NewProgressRepository creates a new progress repository
func NewProgressRepository() *ProgressRepository {
	return &ProgressRepository{
		db: postgres.GetDB(),
	}
}

// Create creates a new progress
func (r *ProgressRepository) Create(progress *models.Progress) error {
	return r.db.Create(progress).Error
}

// GetByID gets a progress by ID
func (r *ProgressRepository) GetByID(id uuid.UUID) (*models.Progress, error) {
	var progress models.Progress
	err := r.db.Preload("Lesson").Where("id = ?", id).First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// GetByUserAndLessonID gets a progress by user ID and lesson ID
func (r *ProgressRepository) GetByUserAndLessonID(userID, lessonID uuid.UUID) (*models.Progress, error) {
	var progress models.Progress
	err := r.db.Where("user_id = ? AND lesson_id = ?", userID, lessonID).First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// Update updates a progress
func (r *ProgressRepository) Update(progress *models.Progress) error {
	return r.db.Save(progress).Error
}

// UpdatePosition updates the position of a progress
func (r *ProgressRepository) UpdatePosition(userID, lessonID uuid.UUID, positionSeconds int) error {
	return r.db.Model(&models.Progress{}).
		Where("user_id = ? AND lesson_id = ?", userID, lessonID).
		Updates(map[string]interface{}{
			"position_seconds": positionSeconds,
			"last_watched_at":  time.Now(),
		}).Error
}

// MarkAsCompleted marks a progress as completed
func (r *ProgressRepository) MarkAsCompleted(userID, lessonID uuid.UUID) error {
	return r.db.Model(&models.Progress{}).
		Where("user_id = ? AND lesson_id = ?", userID, lessonID).
		Updates(map[string]interface{}{
			"completed":       true,
			"last_watched_at": time.Now(),
		}).Error
}

// GetByUserID gets progress by user ID
func (r *ProgressRepository) GetByUserID(userID uuid.UUID) ([]models.Progress, error) {
	var progress []models.Progress
	err := r.db.Preload("Lesson").Where("user_id = ?", userID).Find(&progress).Error
	if err != nil {
		return nil, err
	}
	return progress, nil
}

// GetByCourseAndUserID gets progress by course ID and user ID
func (r *ProgressRepository) GetByCourseAndUserID(courseID, userID uuid.UUID) ([]models.Progress, error) {
	var progress []models.Progress
	err := r.db.Joins("JOIN lessons ON lessons.id = progress.lesson_id").
		Where("lessons.course_id = ? AND progress.user_id = ?", courseID, userID).
		Preload("Lesson").
		Find(&progress).Error
	if err != nil {
		return nil, err
	}
	return progress, nil
}
