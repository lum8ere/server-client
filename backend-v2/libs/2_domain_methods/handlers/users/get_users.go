package users

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GetUsersHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	var users []model.User
	err := sctx.GetDB().Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении состояния объекта: %w", err)
	}

	for _, user := range users {
		user.PasswordHash = ""
	}

	return users, nil
}

func GetUserByIDHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {

	userId, ok := args.GetStringValue("user_id")
	if !ok || userId == "" {
		return nil, fmt.Errorf("user_id is missing")
	}

	var user model.User
	if err := sctx.GetDB().Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	// Не передаём хэш пароля в ответе
	user.PasswordHash = ""
	return user, nil
}

func GetProfileHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	// Извлекаем user_id из ANY_DATA
	userId, ok := args.GetStringValue("user_id")
	if !ok || userId == "" {
		return nil, fmt.Errorf("user_id is missing")
	}

	var user model.User
	if err := sctx.GetDB().Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	// Не передаём хэш пароля в ответе
	user.PasswordHash = ""
	return user, nil
}

// UpdateProfileHandler позволяет изменить данные профиля текущего пользователя.
// Можно обновлять поля username, email и, при необходимости, пароль.
func UpdateProfileHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	userId, ok := args.GetStringValue("user_id")
	if !ok || userId == "" {
		return nil, fmt.Errorf("user_id is missing")
	}

	var user model.User
	if err := sctx.GetDB().Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Обновляем username, если передан
	if username, ok := args.GetStringValue("username"); ok && username != "" {
		user.Username = username
	}

	// Обновляем email, если передан
	if email, ok := args.GetStringValue("email"); ok && email != "" {
		user.Email = email
	}

	// Если передан новый пароль – хэшируем его
	if password, ok := args.GetStringValue("password"); ok && password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hashed)
	}

	user.UpdatedAt = time.Now()

	if err := sctx.GetDB().Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Очищаем поле пароля для ответа
	user.PasswordHash = ""
	return user, nil
}
