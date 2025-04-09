package repositories

import (
	"github.com/0xBoji/web3-edu-core/internal/database/postgres"
	"github.com/0xBoji/web3-edu-core/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: postgres.GetDB(),
	}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

// GetByID gets a refresh token by ID
func (r *RefreshTokenRepository) GetByID(id uuid.UUID) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Where("id = ?", id).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetByToken gets a refresh token by token
func (r *RefreshTokenRepository) GetByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// GetByUserID gets refresh tokens by user ID
func (r *RefreshTokenRepository) GetByUserID(userID uuid.UUID) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	err := r.db.Where("user_id = ?", userID).Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// Delete deletes a refresh token
func (r *RefreshTokenRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.RefreshToken{}, id).Error
}

// DeleteByUserID deletes refresh tokens by user ID
func (r *RefreshTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}
