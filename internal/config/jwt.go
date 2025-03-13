package config

import (
	"github.com/golang-jwt/jwt/v4"
	"sync"
	"time"
)

var (
	tokenBlacklist     = make(map[string]time.Time)
	tokenBlacklistLock sync.RWMutex
)

func GenerateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(LoadConfig().JWTSecret))
}

func ValidateToken(tokenString string) (uint, error) {
	// First, check if token is blacklisted
	tokenBlacklistLock.RLock()
	_, isBlacklisted := tokenBlacklist[tokenString]
	tokenBlacklistLock.RUnlock()

	if isBlacklisted {
		return 0, jwt.ErrSignatureInvalid
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(LoadConfig().JWTSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := uint(claims["user_id"].(float64))
		return userId, nil
	}

	return 0, jwt.ErrSignatureInvalid
}

func BlacklistToken(tokenString string) {
	tokenBlacklistLock.Lock()
	defer tokenBlacklistLock.Unlock()

	// Store token with an expiration time
	tokenBlacklist[tokenString] = time.Now().Add(time.Hour * 24)

	// Periodically clean up expired tokens
	go cleanupBlacklist()
}

func cleanupBlacklist() {
	tokenBlacklistLock.Lock()
	defer tokenBlacklistLock.Unlock()

	now := time.Now()
	for token, expTime := range tokenBlacklist {
		if now.After(expTime) {
			delete(tokenBlacklist, token)
		}
	}
}
