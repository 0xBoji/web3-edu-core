package services

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/0xBoji/web3-edu-core/internal/database/redis"
	"github.com/0xBoji/web3-edu-core/internal/domain/models"
	"github.com/0xBoji/web3-edu-core/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
	cache        *redis.Cache
}

// NewCategoryService creates a new category service
func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repositories.NewCategoryRepository(),
		cache:        redis.NewCache(),
	}
}

// CategoryResponse represents the category response
type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Slug        string    `json:"slug"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCategoryRequest represents the create category request
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

// UpdateCategoryRequest represents the update category request
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

// Create creates a new category
func (s *CategoryService) Create(req CreateCategoryRequest) (*CategoryResponse, error) {
	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	}

	// Check if slug already exists
	_, err := s.categoryRepo.GetBySlug(slug)
	if err == nil {
		return nil, errors.New("slug already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create category
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		Slug:        slug,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.cache.Delete(ctx, "categories:all")

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// GetByID gets a category by ID
func (s *CategoryService) GetByID(id uuid.UUID) (*CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// GetBySlug gets a category by slug
func (s *CategoryService) GetBySlug(slug string) (*CategoryResponse, error) {
	category, err := s.categoryRepo.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// Update updates a category
func (s *CategoryService) Update(id uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Slug != "" {
		// Check if slug already exists
		existingCategory, err := s.categoryRepo.GetBySlug(req.Slug)
		if err == nil && existingCategory.ID != id {
			return nil, errors.New("slug already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			return nil, err
		}
		category.Slug = req.Slug
	}

	// Save category
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.cache.Delete(ctx, "categories:all")

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// Delete deletes a category
func (s *CategoryService) Delete(id uuid.UUID) error {
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	s.cache.Delete(ctx, "categories:all")

	return nil
}

// List lists all categories
func (s *CategoryService) List() ([]CategoryResponse, error) {
	// Try to get from cache
	ctx := context.Background()
	var cachedCategories []CategoryResponse
	err := s.cache.GetJSON(ctx, "categories:all", &cachedCategories)
	if err == nil {
		return cachedCategories, nil
	}

	// Get from database
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, err
	}

	// Convert to response
	var categoryResponses []CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			Slug:        category.Slug,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	// Cache the result
	s.cache.SetJSON(ctx, "categories:all", categoryResponses, 1*time.Hour)

	return categoryResponses, nil
}

// generateSlug generates a slug from a string
func generateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9-]")
	s = reg.ReplaceAllString(s, "")

	// Remove consecutive hyphens
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	return s
}
