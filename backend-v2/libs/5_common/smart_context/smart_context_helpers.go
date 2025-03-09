package smart_context

import (
	"backed-api-v2/libs/5_common/types"
	"sync"
)

const SESSION_ID_KEY = "session_id"

func (sc *SmartContext) WithSessionId(sessionId string) ISmartContext {
	return sc.WithField(SESSION_ID_KEY, sessionId)
}

func (sc *SmartContext) GetSessionId() string {
	result, ok := types.GetFieldTypedValue[string](sc.dataFields, SESSION_ID_KEY)
	if !ok {
		return ""
	}
	return result
}

const WAIT_GROUP = "wait_group"

func (sc *SmartContext) WithWaitGroup(wg *sync.WaitGroup) ISmartContext {
	return sc.WithField(WAIT_GROUP, wg)
}

func (sc *SmartContext) GetWaitGroup() *sync.WaitGroup {
	result, ok := types.GetFieldTypedValue[*sync.WaitGroup](sc.dataFields, WAIT_GROUP)
	if !ok {
		return nil
	}

	return result
}

const DEVICE_ID_KEY = "device_identifier"

func (sc *SmartContext) WithDeviceIdentifier(deviceID string) ISmartContext {
	return sc.WithField(DEVICE_ID_KEY, deviceID)
}

func (sc *SmartContext) GetDeviceIdentifier() string {
	result, ok := types.GetFieldTypedValue[string](sc.dataFields, DEVICE_ID_KEY)
	if !ok {
		return ""
	}
	return result
}

// GEOCODER_KEY – ключ для хранения геокодера в dataFields.
const GEOCODER_KEY = "geocoder"

// WithGeocoder возвращает новый SmartContext с добавленным геокодером.
func (sc *SmartContext) WithGeocoder(geocoderInstance IGeocoder) ISmartContext {
	return sc.WithField(GEOCODER_KEY, geocoderInstance)
}

// GetGeocoder извлекает из SmartContext геокодер, если он был установлен.
func (sc *SmartContext) GetGeocoder() IGeocoder {
	result, ok := types.GetFieldTypedValue[IGeocoder](sc.dataFields, GEOCODER_KEY)
	if !ok {
		return nil
	}
	return result
}
