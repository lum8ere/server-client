package test_handlers

import (
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"fmt"
	"math/rand/v2"
	"net/http"
)

// RndHandler принимает smart_context и дополнительные параметры (ANY_DATA),
// генерирует случайное число и возвращает результат в виде ANY_DATA.
func RndHandler(sctx smart_context.ISmartContext, params types.ANY_DATA, w http.ResponseWriter, r *http.Request) (*types.ANY_DATA, error) {
	sctx.Infof("RndHandler started with params: %v", params)
	randomFloat := rand.Float64()
	result := types.ANY_DATA{
		"result": fmt.Sprintf("%f", randomFloat),
		// Можно добавить возвращаемые параметры или echo входных параметров
		"received_params": params,
	}
	sctx.Infof("RndHandler finished with result: %v", result)
	return types.AnyDataRef(result), nil
}

func RndHandler2(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	sctx.Infof("RndHandler started with args: %v", args)
	randomFloat := rand.Float64()
	result := map[string]any{
		"result":          fmt.Sprintf("%f", randomFloat),
		"received_params": args,
	}
	sctx.Infof("RndHandler finished with result: %v", result)
	return result, nil
}
