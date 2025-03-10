package users

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"fmt"
)

func GetUsersHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	var users []model.User
	err := sctx.GetDB().Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении состояния объекта: %w", err)
	}

	return users, nil
}

func GetUserByIDHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	var user model.User
	err := sctx.GetDB().Where("").First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении состояния объекта: %w", err)
	}

	return user, nil
}
