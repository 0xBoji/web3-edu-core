package handlers

import (
	"net/http"
	"strconv"

	"github.com/0xBoji/web3-edu-core/internal/domain/services"
	"github.com/0xBoji/web3-edu-core/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CourseHandler handles course-related requests
type CourseHandler struct {
	courseService *services.CourseService
}

// NewCourseHandler creates a new course handler
func NewCourseHandler() *CourseHandler {
	return &CourseHandler{
		courseService: services.NewCourseService(),
	}
}

// @Summary Get all courses
// @Description Get a list of all courses with pagination and optional category filter
// @Tags courses
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Param category query string false "Filter by category slug"
// @Success 200 {object} utils.Response{data=[]services.CourseResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses [get]
func (h *CourseHandler) List(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")

	// Get courses
	courses, total, err := h.courseService.List(page, pageSize, category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Page", strconv.Itoa(page))
	c.Header("X-Page-Size", strconv.Itoa(pageSize))

	utils.SuccessResponse(c, courses)
}

// @Summary Get featured courses
// @Description Get a list of featured courses
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]services.CourseResponse}
// @Failure 500 {object} utils.Response
// @Router /courses/featured [get]
func (h *CourseHandler) GetFeatured(c *gin.Context) {
	courses, err := h.courseService.GetFeatured()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, courses)
}

// @Summary Get course by ID
// @Description Get a course by its ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} utils.Response{data=services.CourseResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses/{id} [get]
func (h *CourseHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID")
		return
	}

	course, err := h.courseService.GetByID(id)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, course)
}

// @Summary Get course lessons
// @Description Get lessons for a course
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} utils.Response{data=[]services.LessonBrief}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses/{id}/lessons [get]
func (h *CourseHandler) GetLessons(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID")
		return
	}

	lessons, err := h.courseService.GetLessons(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, lessons)
}

// @Summary Enroll in a course
// @Description Enroll the authenticated user in a course
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses/{id}/enroll [post]
func (h *CourseHandler) Enroll(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse course ID
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID")
		return
	}

	// Enroll user in course
	err = h.courseService.Enroll(userID.(uuid.UUID), courseID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "already enrolled in this course" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

// @Summary Create a course
// @Description Create a new course (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param course body services.CreateCourseRequest true "Course data"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response{data=services.CourseResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/courses [post]
func (h *CourseHandler) Create(c *gin.Context) {
	var req services.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	course, err := h.courseService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, course)
}

// @Summary Update a course
// @Description Update an existing course (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param course body services.UpdateCourseRequest true "Course data"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response{data=services.CourseResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/courses/{id} [put]
func (h *CourseHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID")
		return
	}

	var req services.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	course, err := h.courseService.Update(id, req)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, course)
}

// @Summary Delete a course
// @Description Delete a course (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/courses/{id} [delete]
func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid course ID")
		return
	}

	err = h.courseService.Delete(id)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}
