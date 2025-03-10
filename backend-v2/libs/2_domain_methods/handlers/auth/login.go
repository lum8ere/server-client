package auth

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
)

type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandlerNew выполняет авторизацию пользователя.
// Он извлекает email и password из args, ищет пользователя в БД,
// сравнивает пароль (bcrypt) и генерирует JWT с данными пользователя (включая роль).
func LoginHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	// Извлекаем email
	email, ok := args.GetStringValue("email")
	if !ok || email == "" {
		return nil, fmt.Errorf("missing email")
	}
	// Извлекаем пароль
	password, ok := args.GetStringValue("password")
	if !ok || password == "" {
		return nil, fmt.Errorf("missing password")
	}

	// Поиск пользователя по email
	var user model.User
	if err := sctx.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Сравнение введённого пароля с хэшированным паролем
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Генерация JWT с информацией о пользователе (включая роль)
	token, err := generateJWT(user, sctx, time.Hour*72)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return LoginResponse{Token: token}, nil
}
