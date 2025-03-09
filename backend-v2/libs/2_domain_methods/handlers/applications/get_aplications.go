package applications

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"fmt"
)

func GetApplicationsByDevicesIDHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	// Извлекаем параметр "id" из params
	id, ok := params.GetStringValue("id")
	if !ok || id == "" {
		return nil, fmt.Errorf("missing device id")
	}

	var apps []model.Application
	err := sctx.GetDB().
		Table("applications").
		Joins("JOIN device_applications ON device_applications.application_id = applications.id").
		Where("device_applications.device_id = ?", id).
		Find(&apps).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching applications: %w", err)
	}

	return apps, nil
}
