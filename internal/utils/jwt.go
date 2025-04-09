package utils

import (
	"errors"
	"time"

	"github.com/0xBoji/web3-edu-core/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token
func GenerateToken(userID uuid.UUID, email, role string) (string, error) {
	expireTime := time.Now().Add(time.Duration(config.AppSetting.TokenExpireTime) * time.Hour)
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.AppSetting.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppSetting.JWTSecret))
}

// ParseToken parses the JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppSetting.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken() (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(config.AppSetting.RefreshTokenExpireTime) * time.Hour)
	refreshToken := uuid.New().String()
	return refreshToken, expiresAt, nil
}
