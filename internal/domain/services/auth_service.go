package services

import (
	"context"
	"errors"
	"time"

	"github.com/0xBoji/web3-edu-core/internal/database/redis"
	"github.com/0xBoji/web3-edu-core/internal/domain/models"
	"github.com/0xBoji/web3-edu-core/internal/domain/repositories"
	"github.com/0xBoji/web3-edu-core/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	refreshTokenRepo *repositories.RefreshTokenRepository
	cache            *redis.Cache
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:         repositories.NewUserRepository(),
		refreshTokenRepo: repositories.NewRefreshTokenRepository(),
		cache:            redis.NewCache(),
	}
}

// RegisterRequest represents the register request
type RegisterRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	FullName       string `json:"full_name" binding:"required"`
	Role           string `json:"role"`
	ProfilePicture string `json:"profile_picture"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ForgotPasswordRequest represents the forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the reset password request
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	User         UserResponse `json:"user"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	FullName       string    `json:"full_name"`
	Role           string    `json:"role"`
	ProfilePicture string    `json:"profile_picture,omitempty"`
}

// Register registers a new user
func (s *AuthService) Register(req RegisterRequest) (*TokenResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Create user
	user := &models.User{
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		FullName:       req.FullName,
		Role:           req.Role,
		ProfilePicture: req.ProfilePicture,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateTokens(user)
}

// Login logs in a user
func (s *AuthService) Login(req LoginRequest) (*TokenResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	return s.generateTokens(user)
}

// RefreshToken refreshes the access token
func (s *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error) {
	// Get refresh token
	token, err := s.refreshTokenRepo.GetByToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if token is expired
	if token.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(token.UserID)
	if err != nil {
		return nil, err
	}

	// Delete old refresh token
	if err := s.refreshTokenRepo.Delete(token.ID); err != nil {
		return nil, err
	}

	// Generate new tokens
	return s.generateTokens(user)
}

// Logout logs out a user
func (s *AuthService) Logout(refreshToken string) error {
	// Get refresh token
	token, err := s.refreshTokenRepo.GetByToken(refreshToken)
	if err != nil {
		return nil // Ignore error if token not found
	}

	// Delete refresh token
	return s.refreshTokenRepo.Delete(token.ID)
}

// ForgotPassword initiates the forgot password process
func (s *AuthService) ForgotPassword(req ForgotPasswordRequest) error {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // Don't reveal that the email doesn't exist
		}
		return err
	}

	// Generate reset token
	token := uuid.New().String()

	// Store token in Redis with expiration
	ctx := context.Background()
	key := "reset_token:" + token
	userData := map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
	}
	if err := s.cache.SetJSON(ctx, key, userData, 1*time.Hour); err != nil {
		return err
	}

	// In a real application, you would send an email with the reset link
	// For now, we'll just return success
	return nil
}

// ResetPassword resets a user's password
func (s *AuthService) ResetPassword(req ResetPasswordRequest) error {
	// Get token from Redis
	ctx := context.Background()
	key := "reset_token:" + req.Token
	var userData map[string]string
	if err := s.cache.GetJSON(ctx, key, &userData); err != nil {
		return errors.New("invalid or expired token")
	}

	// Parse user ID
	userID, err := uuid.Parse(userData["user_id"])
	if err != nil {
		return err
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Update user password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Delete token from Redis
	if err := s.cache.Delete(ctx, key); err != nil {
		return err
	}

	// Delete all refresh tokens for this user
	return s.refreshTokenRepo.DeleteByUserID(user.ID)
}

// generateTokens generates access and refresh tokens
func (s *AuthService) generateTokens(user *models.User) (*TokenResponse, error) {
	// Generate access token
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, expiresAt, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Save refresh token
	token := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}

	if err := s.refreshTokenRepo.Create(token); err != nil {
		return nil, err
	}

	// Create response
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: UserResponse{
			ID:             user.ID,
			Email:          user.Email,
			FullName:       user.FullName,
			Role:           user.Role,
			ProfilePicture: user.ProfilePicture,
		},
	}, nil
}
