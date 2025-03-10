package auth

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterResponse – структура ответа на регистрацию
type RegisterResponse struct {
	Token string `json:"token"`
}

// RegisterHandler регистрирует нового пользователя с назначением роли по коду приглашения.
// Используем формат входных данных через types.ANY_DATA.
func RegisterHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	username, ok := args.GetStringValue("username")
	if !ok || username == "" {
		return nil, fmt.Errorf("missing username")
	}
	email, ok := args.GetStringValue("email")
	if !ok || email == "" {
		return nil, fmt.Errorf("missing email")
	}
	password, ok := args.GetStringValue("password")
	if !ok || password == "" {
		return nil, fmt.Errorf("missing password")
	}
	// invitationCode, _ := args.GetStringValue("invitation_code")

	// Хэширование пароля с использованием bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		sctx.Errorf("Error hashing password: %v", err)
		return nil, fmt.Errorf("error processing request")
	}

	// Определяем роль – по умолчанию "observer"
	role := "OBSERVER"
	// if invitationCode != "" {
	// 	switch invitationCode {
	// 	case "admininvite":
	// 		role = "admin"
	// 	case "observerplusinvite":
	// 		role = "observer+"
	// 	default:
	// 		return nil, fmt.Errorf("invalid invitation code")
	// 	}
	// 	// Здесь можно отметить использование приглашения (например, сохранить в таблице приглашений)
	// }

	// Создаем нового пользователя
	newUser := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleCode:     role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := sctx.GetDB().Create(newUser).Error; err != nil {
		sctx.Errorf("Error creating user: %v", err)
		return nil, fmt.Errorf("error creating user")
	}

	// Генерация JWT-токена с данными пользователя (включая роль)
	token, err := generateJWT(*newUser, sctx, time.Hour*72)
	if err != nil {
		sctx.Errorf("Error generating JWT: %v", err)
		return nil, fmt.Errorf("error generating token")
	}

	return RegisterResponse{Token: token}, nil
}
