package devices

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RegisterDeviceRequest – структура запроса на регистрацию устройства.
type RegisterDeviceRequest struct {
	DeviceIdentifier string `json:"device_identifier"`
	Description      string `json:"description,omitempty"`
	// При необходимости можно добавить дополнительные поля, например, UserID
}

func RegisterDeviceHandler(sctx smart_context.ISmartContext, params types.ANY_DATA, w http.ResponseWriter, r *http.Request) (*types.ANY_DATA, error) {
	sctx.Infof("RndHandler started with params: %v", params)
	var req RegisterDeviceRequest
	paramsData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %v", err)
	}
	err = json.Unmarshal(paramsData, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal object: %v", err)
	}

	device := model.Device{
		DeviceIdentifier: req.DeviceIdentifier,
		Description:      req.Description,
		Status:           "ONLINE", // можно использовать значение статуса из базы (например, "online")
		LastSeen:         time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = sctx.GetDB().Create(&device).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении состояния объекта: %w", err)
	}

	return nil, nil
}