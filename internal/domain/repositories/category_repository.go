package repositories

import (
	"github.com/0xBoji/web3-edu-core/internal/database/postgres"
	"github.com/0xBoji/web3-edu-core/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		db: postgres.GetDB(),
	}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// GetByID gets a category by ID
func (r *CategoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetBySlug gets a category by slug
func (r *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// List lists all categories
func (r *CategoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCourses gets courses by category ID
func (r *CategoryRepository) GetCourses(categoryID uuid.UUID, page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var count int64

	query := r.db.Joins("JOIN course_categories ON courses.id = course_categories.course_id").
		Where("course_categories.category_id = ?", categoryID)

	query.Count(&count)

	offset := (page - 1) * pageSize
	err := query.Preload("Instructor").Offset(offset).Limit(pageSize).Find(&courses).Error
	if err != nil {
		return nil, 0, err
	}

	return courses, count, nil
}
