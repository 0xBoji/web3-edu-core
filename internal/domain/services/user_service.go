package services

import (
	"errors"

	"github.com/0xBoji/web3-edu-core/internal/domain/repositories"
	"github.com/0xBoji/web3-edu-core/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
	}
}

// We'll use the UserResponse from auth_service.go

// UpdateUserRequest represents the update user request
type UpdateUserRequest struct {
	FullName       string `json:"full_name"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture"`
}

// GetByID gets a user by ID
func (s *UserService) GetByID(id uuid.UUID) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FullName:       user.FullName,
		Role:           user.Role,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

// Update updates a user
func (s *UserService) Update(id uuid.UUID, req UpdateUserRequest) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.ProfilePicture != "" {
		user.ProfilePicture = req.ProfilePicture
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashedPassword
	}

	// Save user
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FullName:       user.FullName,
		Role:           user.Role,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

// List lists all users
func (s *UserService) List(page, pageSize int) ([]UserResponse, int64, error) {
	users, count, err := s.userRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:             user.ID,
			Email:          user.Email,
			FullName:       user.FullName,
			Role:           user.Role,
			ProfilePicture: user.ProfilePicture,
		})
	}

	return userResponses, count, nil
}

// Delete deletes a user
func (s *UserService) Delete(id uuid.UUID) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.Delete(id)
}
