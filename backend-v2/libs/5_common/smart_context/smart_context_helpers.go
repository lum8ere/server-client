package smart_context

import "backed-api-v2/libs/5_common/types"

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