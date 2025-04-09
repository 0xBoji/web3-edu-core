package handlers

import (
	"github.com/0xBoji/web3-edu-core/internal/domain/services"
	"github.com/0xBoji/web3-edu-core/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: services.NewCategoryService(),
	}
}

// List handles the list categories request
// @Summary List all categories
// @Description Get a list of all categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]services.CategoryResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Router /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.List()
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, categories)
}

// Get handles the get category request
// @Summary Get a category by ID
// @Description Get a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} utils.Response{data=services.CategoryResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Router /categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid category ID")
		return
	}

	category, err := h.categoryService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, category)
}

// Create handles the create category request
// @Summary Create a new category
// @Description Create a new category (admin only)
// @Tags admin,categories
// @Accept json
// @Produce json
// @Param request body services.CreateCategoryRequest true "Create Category Request"
// @Success 200 {object} utils.Response{data=services.CategoryResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Security BearerAuth
// @Router /admin/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req services.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	category, err := h.categoryService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, category)
}

// Update handles the update category request
// @Summary Update a category
// @Description Update a category (admin only)
// @Tags admin,categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body services.UpdateCategoryRequest true "Update Category Request"
// @Success 200 {object} utils.Response{data=services.CategoryResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Security BearerAuth
// @Router /admin/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid category ID")
		return
	}

	var req services.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	category, err := h.categoryService.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, category)
}

// Delete handles the delete category request
// @Summary Delete a category
// @Description Delete a category (admin only)
// @Tags admin,categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} utils.Response{data=object{message=string}} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Security BearerAuth
// @Router /admin/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid category ID")
		return
	}

	if err := h.categoryService.Delete(id); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "category deleted successfully"})
}
