package repositories

import (
	"github.com/0xBoji/web3-edu-core/internal/database/postgres"
	"github.com/0xBoji/web3-edu-core/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

// NewEnrollmentRepository creates a new enrollment repository
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		db: postgres.GetDB(),
	}
}

// Create creates a new enrollment
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

// GetByID gets an enrollment by ID
func (r *EnrollmentRepository) GetByID(id uuid.UUID) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Preload("User").Preload("Course").Where("id = ?", id).First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// Delete deletes an enrollment
func (r *EnrollmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Enrollment{}, id).Error
}

// GetByUserID gets enrollments by user ID
func (r *EnrollmentRepository) GetByUserID(userID uuid.UUID, page, pageSize int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var count int64

	r.db.Model(&models.Enrollment{}).Where("user_id = ?", userID).Count(&count)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Course").Preload("Course.Instructor").Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&enrollments).Error
	if err != nil {
		return nil, 0, err
	}

	return enrollments, count, nil
}

// GetByCourseID gets enrollments by course ID
func (r *EnrollmentRepository) GetByCourseID(courseID uuid.UUID, page, pageSize int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var count int64

	r.db.Model(&models.Enrollment{}).Where("course_id = ?", courseID).Count(&count)

	offset := (page - 1) * pageSize
	err := r.db.Preload("User").Where("course_id = ?", courseID).Offset(offset).Limit(pageSize).Find(&enrollments).Error
	if err != nil {
		return nil, 0, err
	}

	return enrollments, count, nil
}

// IsEnrolled checks if a user is enrolled in a course
func (r *EnrollmentRepository) IsEnrolled(userID, courseID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Enrollment{}).Where("user_id = ? AND course_id = ?", userID, courseID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
