package dicts

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"fmt"
)

func GetRoleDictsHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	var roles []model.Role
	err := sctx.GetDB().Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении состояния объекта: %w", err)
	}

	return roles, nil
}

func GetGroupsDictsHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	var roles []model.DeviceGroup
	err := sctx.GetDB().Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении состояния объекта: %w", err)
	}

	return roles, nil
}
