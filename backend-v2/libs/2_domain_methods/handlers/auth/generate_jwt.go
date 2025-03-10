package auth

import (
	"fmt"
	"os"
	"time"

	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"

	"github.com/golang-jwt/jwt/v4"
)

// generateJWT генерирует JWT-токен, включающий данные пользователя (user_id, username, role) и срок действия.
func generateJWT(user model.User, sctx smart_context.ISmartContext, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.RoleCode,
		"exp":      time.Now().Add(duration).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}
	return token.SignedString([]byte(secret))
}
