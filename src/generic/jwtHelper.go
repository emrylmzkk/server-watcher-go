package generic

import (
	"errors"
	"os"
	"time"

	"server-watcher-app/src/models"
	enumModels "server-watcher-app/src/models/enum"

	"github.com/golang-jwt/jwt/v5"
)

// Claims yapısına UserRole ekledik, böylece Middleware'lerde rol kontrolü kolaylaşır.
type Claims struct {
	UserID   uint                `json:"user_id"`
	Username string              `json:"username"`
	Role     enumModels.UserRole `json:"role"`
	Type     string              `json:"type"` // "access" | "refresh"
	jwt.RegisteredClaims
}

// GenerateTokenPair artık direkt modeldeki User nesnesini alıyor.
func GenerateTokenPair(user *models.User) (string, string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	// Access Token - 15 Dakika
	access, err := generateToken(user, "access", 15*time.Minute, secret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token - 7 Gün
	refresh, err := generateToken(user, "refresh", 7*24*time.Hour, secret)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func ValidateToken(tokenStr string) (*Claims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func generateToken(user *models.User, tokenType string, duration time.Duration, secret []byte) (string, error) {
	claims := &Claims{
		UserID:   user.ID, // gorm.Model'den gelen uint ID
		Username: user.Username,
		Role:     user.UserRole,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "server-watcher-app",
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}
