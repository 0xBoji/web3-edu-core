package handlers

import (
	"github.com/0xBoji/web3-edu-core/internal/domain/services"
	"github.com/0xBoji/web3-edu-core/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// Register handles the register request
// @Summary Register a new user
// @Description Register a new user with email, password, and full name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.RegisterRequest true "Register Request"
// @Success 200 {object} utils.Response{data=services.TokenResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	resp, err := h.authService.Register(req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, resp)
}

// Login handles the login request
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login Request"
// @Success 200 {object} utils.Response{data=services.TokenResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, resp)
}

// RefreshToken handles the refresh token request
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "Refresh Token Request"
// @Success 200 {object} utils.Response{data=services.TokenResponse} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, resp)
}

// Logout handles the logout request
// @Summary Logout a user
// @Description Logout a user by invalidating refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "Logout Request"
// @Success 200 {object} utils.Response{data=object{message=string}} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "logged out successfully"})
}

// ForgotPassword handles the forgot password request
// @Summary Forgot password
// @Description Initiate the forgot password process
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} utils.Response{data=object{message=string}} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req services.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.authService.ForgotPassword(req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Always return success to prevent email enumeration
	utils.SuccessResponse(c, gin.H{"message": "if your email is registered, you will receive a password reset link"})
}

// ResetPassword handles the reset password request
// @Summary Reset password
// @Description Reset password using token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} utils.Response{data=object{message=string}} "Success"
// @Failure 400 {object} utils.Response "Bad Request"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req services.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.authService.ResetPassword(req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "password reset successfully"})
}
