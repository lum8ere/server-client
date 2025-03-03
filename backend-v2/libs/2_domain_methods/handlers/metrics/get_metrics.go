package metrics

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"fmt"
)

func GetMetricsHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	var metrics []model.Metric
	err := sctx.GetDB().First(&metrics).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении состояния объекта: %w", err)
	}

	return metrics, nil
}

func GetMetricsByDeviceIDHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	// Извлекаем параметр "id" из params
	id, ok := params.GetStringValue("id")
	if !ok || id == "" {
		return nil, fmt.Errorf("missing device id")
	}

	var metrics model.Metric
	err := sctx.GetDB().Where("device_id = ?", id).First(&metrics).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении устройства: %w", err)
	}

	return metrics, nil
}
