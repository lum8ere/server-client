package devices

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"fmt"
)

func GetDevicesHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	var devices []model.Device
	err := sctx.GetDB().Find(&devices).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении состояния объекта: %w", err)
	}

	return devices, nil
}

func GetDevicesByIDHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	// Извлекаем параметр "id" из params
	id, ok := params.GetStringValue("id")
	if !ok || id == "" {
		return nil, fmt.Errorf("missing device id")
	}

	var device model.Device
	// Ищем устройство по id
	err := sctx.GetDB().Where("id = ?", id).First(&device).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении устройства: %w", err)
	}

	return device, nil
}
