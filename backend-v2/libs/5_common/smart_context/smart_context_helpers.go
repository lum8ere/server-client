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
