package handlers

import (
	"strconv"

	"github.com/0xBoji/web3-edu-core/internal/domain/services"
	"github.com/0xBoji/web3-edu-core/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// Get handles the get user request
// @Summary Get a user by ID
// @Description Get a user by ID (admin can view any user, regular users can only view themselves)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response{data=services.UserResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid user ID")
		return
	}

	// Check if the user is requesting their own profile or is an admin
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if id != userID.(uuid.UUID) && role.(string) != "admin" {
		utils.ForbiddenResponse(c)
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

// Update handles the update user request
// @Summary Update a user
// @Description Update a user (admin can update any user, regular users can only update themselves)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body services.UpdateUserRequest true "Update User Request"
// @Success 200 {object} utils.Response{data=services.UserResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid user ID")
		return
	}

	// Check if the user is updating their own profile or is an admin
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if id != userID.(uuid.UUID) && role.(string) != "admin" {
		utils.ForbiddenResponse(c)
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.userService.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

// List handles the list users request
// @Summary List all users
// @Description List all users (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} utils.Response{data=object{users=[]services.UserResponse,total=int,page=int,size=int}} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	// Only admin can list all users
	role, _ := c.Get("role")
	if role.(string) != "admin" {
		utils.ForbiddenResponse(c)
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, count, err := h.userService.List(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"users": users,
		"total": count,
		"page":  page,
		"size":  pageSize,
	})
}

// Delete handles the delete user request
// @Summary Delete a user
// @Description Delete a user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response{data=object{message=string}} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid user ID")
		return
	}

	// Only admin can delete users
	role, _ := c.Get("role")
	if role.(string) != "admin" {
		utils.ForbiddenResponse(c)
		return
	}

	if err := h.userService.Delete(id); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "user deleted successfully"})
}

// GetProfile handles the get profile request
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=services.UserResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Security BearerAuth
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(uuid.UUID)

	user, err := h.userService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

// UpdateProfile handles the update profile request
// @Summary Update current user profile
// @Description Update the profile of the currently authenticated user
// @Tags profile
// @Accept json
// @Produce json
// @Param request body services.UpdateUserRequest true "Update User Request"
// @Success 200 {object} utils.Response{data=services.UserResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Security BearerAuth
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(uuid.UUID)

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.userService.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}
