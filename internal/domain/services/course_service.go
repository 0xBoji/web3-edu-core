package services

import (
	"context"
	"errors"
	"time"

	"github.com/0xBoji/web3-edu-core/internal/database/redis"
	"github.com/0xBoji/web3-edu-core/internal/domain/models"
	"github.com/0xBoji/web3-edu-core/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseService struct {
	courseRepo     *repositories.CourseRepository
	lessonRepo     *repositories.LessonRepository
	enrollmentRepo *repositories.EnrollmentRepository
	cache          *redis.Cache
}

// NewCourseService creates a new course service
func NewCourseService() *CourseService {
	return &CourseService{
		courseRepo:     repositories.NewCourseRepository(),
		lessonRepo:     repositories.NewLessonRepository(),
		enrollmentRepo: repositories.NewEnrollmentRepository(),
		cache:          redis.NewCache(),
	}
}

// CourseResponse represents the course response
type CourseResponse struct {
	ID           uuid.UUID     `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description,omitempty"`
	Thumbnail    string        `json:"thumbnail,omitempty"`
	InstructorID uuid.UUID     `json:"instructor_id"`
	Instructor   UserResponse  `json:"instructor,omitempty"`
	Price        float64       `json:"price"`
	Level        string        `json:"level,omitempty"`
	Duration     int           `json:"duration,omitempty"`
	Category     string        `json:"category,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Lessons      []LessonBrief `json:"lessons,omitempty"`
}

// LessonBrief represents a brief version of a lesson
type LessonBrief struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Duration    int       `json:"duration,omitempty"`
	OrderNumber int       `json:"order_number"`
}

// CreateCourseRequest represents the create course request
type CreateCourseRequest struct {
	Title        string    `json:"title" binding:"required"`
	Description  string    `json:"description"`
	Thumbnail    string    `json:"thumbnail"`
	InstructorID uuid.UUID `json:"instructor_id" binding:"required"`
	Price        float64   `json:"price"`
	Level        string    `json:"level"`
	Duration     int       `json:"duration"`
	Category     string    `json:"category"`
}

// UpdateCourseRequest represents the update course request
type UpdateCourseRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Thumbnail   string  `json:"thumbnail"`
	Price       float64 `json:"price"`
	Level       string  `json:"level"`
	Duration    int     `json:"duration"`
	Category    string  `json:"category"`
}

// Create creates a new course
func (s *CourseService) Create(req CreateCourseRequest) (*CourseResponse, error) {
	course := &models.Course{
		Title:        req.Title,
		Description:  req.Description,
		Thumbnail:    req.Thumbnail,
		InstructorID: req.InstructorID,
		Price:        req.Price,
		Level:        req.Level,
		Duration:     req.Duration,
		Category:     req.Category,
	}

	if err := s.courseRepo.Create(course); err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.cache.Delete(ctx, "courses:list")
	s.cache.Delete(ctx, "courses:featured")

	return s.mapCourseToResponse(course), nil
}

// GetByID gets a course by ID
func (s *CourseService) GetByID(id uuid.UUID) (*CourseResponse, error) {
	// Try to get from cache
	ctx := context.Background()
	var cachedCourse CourseResponse
	err := s.cache.GetJSON(ctx, "course:"+id.String(), &cachedCourse)
	if err == nil {
		return &cachedCourse, nil
	}

	// Get from database
	course, err := s.courseRepo.GetByIDWithLessons(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	response := s.mapCourseToResponse(course)

	// Cache the result
	s.cache.SetJSON(ctx, "course:"+id.String(), response, 1*time.Hour)

	return response, nil
}

// Update updates a course
func (s *CourseService) Update(id uuid.UUID, req UpdateCourseRequest) (*CourseResponse, error) {
	course, err := s.courseRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	// Update fields
	if req.Title != "" {
		course.Title = req.Title
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.Thumbnail != "" {
		course.Thumbnail = req.Thumbnail
	}
	if req.Price != 0 {
		course.Price = req.Price
	}
	if req.Level != "" {
		course.Level = req.Level
	}
	if req.Duration != 0 {
		course.Duration = req.Duration
	}
	if req.Category != "" {
		course.Category = req.Category
	}

	// Save course
	if err := s.courseRepo.Update(course); err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.cache.Delete(ctx, "course:"+id.String())
	s.cache.Delete(ctx, "courses:list")
	s.cache.Delete(ctx, "courses:featured")

	return s.mapCourseToResponse(course), nil
}

// Delete deletes a course
func (s *CourseService) Delete(id uuid.UUID) error {
	_, err := s.courseRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return err
	}

	if err := s.courseRepo.Delete(id); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	s.cache.Delete(ctx, "course:"+id.String())
	s.cache.Delete(ctx, "courses:list")
	s.cache.Delete(ctx, "courses:featured")

	return nil
}

// List lists all courses
func (s *CourseService) List(page, pageSize int, category string) ([]CourseResponse, int64, error) {
	// Try to get from cache if no filters
	if category == "" {
		ctx := context.Background()
		var cachedCourses struct {
			Courses []CourseResponse `json:"courses"`
			Total   int64            `json:"total"`
		}
		cacheKey := "courses:list:page:" + string(rune(page)) + ":size:" + string(rune(pageSize))
		err := s.cache.GetJSON(ctx, cacheKey, &cachedCourses)
		if err == nil {
			return cachedCourses.Courses, cachedCourses.Total, nil
		}
	}

	// Get from database
	var courses []models.Course
	var count int64
	var err error

	if category != "" {
		courses, count, err = s.courseRepo.ListByCategory(category, page, pageSize)
	} else {
		courses, count, err = s.courseRepo.List(page, pageSize)
	}

	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	var courseResponses []CourseResponse
	for _, course := range courses {
		courseResponses = append(courseResponses, *s.mapCourseToResponse(&course))
	}

	// Cache the result if no filters
	if category == "" {
		ctx := context.Background()
		cacheKey := "courses:list:page:" + string(rune(page)) + ":size:" + string(rune(pageSize))
		s.cache.SetJSON(ctx, cacheKey, struct {
			Courses []CourseResponse `json:"courses"`
			Total   int64            `json:"total"`
		}{
			Courses: courseResponses,
			Total:   count,
		}, 30*time.Minute)
	}

	return courseResponses, count, nil
}

// GetFeatured gets featured courses
func (s *CourseService) GetFeatured() ([]CourseResponse, error) {
	// Try to get from cache
	ctx := context.Background()
	var cachedCourses []CourseResponse
	err := s.cache.GetJSON(ctx, "courses:featured", &cachedCourses)
	if err == nil {
		return cachedCourses, nil
	}

	// Get from database (for now, just return the first 5 courses)
	courses, _, err := s.courseRepo.List(1, 5)
	if err != nil {
		return nil, err
	}

	// Convert to response
	var courseResponses []CourseResponse
	for _, course := range courses {
		courseResponses = append(courseResponses, *s.mapCourseToResponse(&course))
	}

	// Cache the result
	s.cache.SetJSON(ctx, "courses:featured", courseResponses, 1*time.Hour)

	return courseResponses, nil
}

// GetLessons gets lessons for a course
func (s *CourseService) GetLessons(courseID uuid.UUID) ([]LessonBrief, error) {
	lessons, err := s.lessonRepo.GetByCourseID(courseID)
	if err != nil {
		return nil, err
	}

	var lessonResponses []LessonBrief
	for _, lesson := range lessons {
		lessonResponses = append(lessonResponses, LessonBrief{
			ID:          lesson.ID,
			Title:       lesson.Title,
			Description: lesson.Description,
			Duration:    lesson.Duration,
			OrderNumber: lesson.OrderNumber,
		})
	}

	return lessonResponses, nil
}

// Enroll enrolls a user in a course
func (s *CourseService) Enroll(userID, courseID uuid.UUID) error {
	// Check if already enrolled
	enrolled, err := s.enrollmentRepo.IsEnrolled(userID, courseID)
	if err != nil {
		return err
	}
	if enrolled {
		return errors.New("already enrolled in this course")
	}

	// Check if course exists
	_, err = s.courseRepo.GetByID(courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return err
	}

	// Create enrollment
	enrollment := &models.Enrollment{
		UserID:   userID,
		CourseID: courseID,
	}

	return s.enrollmentRepo.Create(enrollment)
}

// mapCourseToResponse maps a course model to a course response
func (s *CourseService) mapCourseToResponse(course *models.Course) *CourseResponse {
	response := &CourseResponse{
		ID:           course.ID,
		Title:        course.Title,
		Description:  course.Description,
		Thumbnail:    course.Thumbnail,
		InstructorID: course.InstructorID,
		Price:        course.Price,
		Level:        course.Level,
		Duration:     course.Duration,
		Category:     course.Category,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
	}

	if course.Instructor.ID != uuid.Nil {
		response.Instructor = UserResponse{
			ID:             course.Instructor.ID,
			Email:          course.Instructor.Email,
			FullName:       course.Instructor.FullName,
			Role:           course.Instructor.Role,
			ProfilePicture: course.Instructor.ProfilePicture,
		}
	}

	if course.Lessons != nil {
		for _, lesson := range course.Lessons {
			response.Lessons = append(response.Lessons, LessonBrief{
				ID:          lesson.ID,
				Title:       lesson.Title,
				Description: lesson.Description,
				Duration:    lesson.Duration,
				OrderNumber: lesson.OrderNumber,
			})
		}
	}

	return response
}
