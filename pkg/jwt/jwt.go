package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/agnathor/finances-go/internal/config"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret           string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewManager(cfg config.JWTConfig) *Manager {
	return &Manager{
		secret:           cfg.Secret,
		accessExpiration:  cfg.AccessExpiration,
		refreshExpiration: cfg.RefreshExpiration,
	}
}

func (m *Manager) GenerateAccessToken(userID string) (string, int64, error) {
	expiresAt := time.Now().Add(m.accessExpiration)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "finances-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, expiresAt.Unix(), nil
}

func (m *Manager) GenerateRefreshToken(userID string) (string, string, int64, error) {
	expiresAt := time.Now().Add(m.refreshExpiration)

	tokenID := fmt.Sprintf("ref_%s_%d", userID, time.Now().UnixNano())

	claims := jwt.MapClaims{
		"user_id":  userID,
		"token_id": tokenID,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "finances-api",
		"type":     "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return signed, tokenID, expiresAt.Unix(), nil
}

func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (m *Manager) ValidateRefreshToken(tokenString string) (userID, tokenID string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("invalid refresh token claims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", fmt.Errorf("invalid token type")
	}

	userID, ok = claims["user_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid user_id in token")
	}

	tokenID, ok = claims["token_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid token_id in token")
	}

	return userID, tokenID, nil
}

func (m *Manager) AccessExpiration() time.Duration {
	return m.accessExpiration
}
