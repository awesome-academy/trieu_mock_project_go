package utils

import (
	"time"
	"trieu_mock_project_go/internal/config"
	appErrors "trieu_mock_project_go/internal/errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWTToken generates a JWT token with user claims
func GenerateJWTToken(userID uint, email string) (string, error) {
	cfg := config.LoadConfig()

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWTToken parses and validates JWT token, returns claims
func ParseJWTToken(tokenString string) (*JWTClaims, error) {
	cfg := config.LoadConfig()

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appErrors.ErrUnexpectedSigningMethod
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, appErrors.ErrInvalidToken
	}

	if !token.Valid {
		return nil, appErrors.ErrInvalidToken
	}

	return claims, nil
}
