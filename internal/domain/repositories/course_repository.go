package repositories

import (
	"github.com/0xBoji/web3-edu-core/internal/database/postgres"
	"github.com/0xBoji/web3-edu-core/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new course repository
func NewCourseRepository() *CourseRepository {
	return &CourseRepository{
		db: postgres.GetDB(),
	}
}

// Create creates a new course
func (r *CourseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

// GetByID gets a course by ID
func (r *CourseRepository) GetByID(id uuid.UUID) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Instructor").Where("id = ?", id).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// GetByIDWithLessons gets a course by ID with lessons
func (r *CourseRepository) GetByIDWithLessons(id uuid.UUID) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Instructor").Preload("Lessons", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_number ASC")
	}).Where("id = ?", id).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// Update updates a course
func (r *CourseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

// Delete deletes a course
func (r *CourseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Course{}, id).Error
}

// List lists all courses
func (r *CourseRepository) List(page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var count int64

	r.db.Model(&models.Course{}).Count(&count)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Instructor").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&courses).Error
	if err != nil {
		return nil, 0, err
	}

	return courses, count, nil
}

// ListByCategory lists courses by category with pagination
func (r *CourseRepository) ListByCategory(category string, page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var count int64

	r.db.Model(&models.Course{}).Where("category = ?", category).Count(&count)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Instructor").Where("category = ?", category).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&courses).Error
	if err != nil {
		return nil, 0, err
	}

	return courses, count, nil
}

// GetByInstructorID gets courses by instructor ID
func (r *CourseRepository) GetByInstructorID(instructorID uuid.UUID, page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var count int64

	r.db.Model(&models.Course{}).Where("instructor_id = ?", instructorID).Count(&count)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Instructor").Where("instructor_id = ?", instructorID).Offset(offset).Limit(pageSize).Find(&courses).Error
	if err != nil {
		return nil, 0, err
	}

	return courses, count, nil
}
