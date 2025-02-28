package auth

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"backed-api-v2/libs/5_common/smart_context"
)

// Здесь определим константы для провайдеров.
const (
	ProviderUnknown          = "UNKNOWN"
)

// Предположим, что у нас есть модель пользователя и credential в базе данных.
type User struct {
	ID           uint
	Username     string
	PasswordHash string
	IsActive     bool
	Role         string
}

// fullLogin – упрощённая логика полного логина, которая проверяет учетные данные,
// генерирует access и refresh токены и сохраняет сессию (например, в Redis).
func fullLogin(sctx smart_context.ISmartContext, req LoginRequest) (*LoginResponse, error) {
	// Здесь вы переключаете контекст на tenant-специфичную БД, если требуется.
	db := sctx.GetDB()

	var user User
	// Выполняем поиск пользователя по логину (например, по email или username, приведенному к нижнему регистру).
	if err := db.Where("LOWER(username) = LOWER(?)", req.Login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Здесь должна быть проверка пароля через bcrypt (пример опущен для краткости).
	// Если пароль не совпадает:
	// return nil, errors.New("invalid credentials")

	// Если пользователь не активен, возвращаем ошибку.
	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	// Генерируем access и refresh токены.
	accessToken, err := generateJWT(user, sctx, time.Hour*72)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := generateJWT(user, sctx, time.Hour*168) // например, на неделю
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Сохраняем сессию в Redis, если требуется (логика не приведена).
	// Например: saveSessionInRedis(sctx, user.ID, refreshToken)

	resp := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return resp, nil
}

// Login выбирает схему аутентификации по провайдеру.
func Login(sctx smart_context.ISmartContext, req LoginRequest) (*LoginResponse, error) {
	// Для примера реализуем только полное логинирование.
	switch req.ProviderCode {
	case ProviderUnknown:
		return fullLogin(sctx, req)
	default:
		return nil, fmt.Errorf("unsupported provider_code: %s", req.ProviderCode)
	}
}

// generateJWT генерирует JWT-токен на основе данных пользователя.
// Здесь мы используем secret из smart_context или из переменных окружения.
func generateJWT(user User, sctx smart_context.ISmartContext, duration time.Duration) (string, error) {
	// Пример генерации токена (используйте библиотеку github.com/golang-jwt/jwt)
	// claims := jwt.MapClaims{
	// 	"user_id":  user.ID,
	// 	"username": user.Username,
	// 	"role":     user.Role,
	// 	"exp":      time.Now().Add(duration).Unix(),
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// secret := []byte(os.Getenv("JWT_SECRET"))
	// return token.SignedString(secret)
	// Для примера вернем заглушку:
	return "dummy.jwt.token", nil
}
