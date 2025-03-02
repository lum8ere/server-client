package run_processor

import (
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"encoding/json"
	"net/http"
	"strings"
)

// HandlerFuncWithReturnAndParams определяет функцию-хендлер, которая принимает smart_context,
// дополнительные параметры в виде types.ANY_DATA, а также стандартные аргументы HTTP.
// Возвращает указатель на types.ANY_DATA и ошибку.
type HandlerFuncWithReturnAndParams func(sctx smart_context.ISmartContext, params types.ANY_DATA, w http.ResponseWriter, r *http.Request) (*types.ANY_DATA, error)
// SmartHandlerFunc принимает smart-контекст и параметры запроса, возвращает объект (любой тип) и ошибку. Если объект не nil, он будет сериализован в JSON.
type SmartHandlerFunc func(sctx smart_context.ISmartContext, args types.ANY_DATA) (*types.ANY_DATA, error)

func WrapHandlerWithReturnAndParams(sctx smart_context.ISmartContext, handler HandlerFuncWithReturnAndParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем параметры из query-параметров:
		params := types.ANY_DATA{}
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}

		// Если метод POST и Content-Type равен application/json, пытаемся декодировать тело запроса
		if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
			decoder := json.NewDecoder(r.Body)
			// Если в теле пришел JSON-объект, то объединяем его с уже полученными query-параметрами
			var bodyParams types.ANY_DATA
			if err := decoder.Decode(&bodyParams); err == nil {
				for k, v := range bodyParams {
					params[k] = v
				}
			}
		}

		result, err := handler(sctx, params, w, r)
		if err != nil {
			sctx.Errorf("Handler error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			sctx.Errorf("Error encoding response: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
	}
}

func WrapSmartHandler(sctx smart_context.ISmartContext, handler SmartHandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Извлечение query-параметров
        params := types.ANY_DATA{}
        for key, values := range r.URL.Query() {
            if len(values) > 0 {
                params[key] = values[0]
            }
        }

        // Если метод POST и Content-Type содержит "application/json", пытаемся декодировать тело запроса
        if r.Method == http.MethodPost && strings.Contains(r.Header.Get("Content-Type"), "application/json") {
            var bodyParams types.ANY_DATA
            if err := json.NewDecoder(r.Body).Decode(&bodyParams); err == nil {
                for k, v := range bodyParams {
                    params[k] = v
                }
            }
        }

        // Вызов основного хендлера с переданными параметрами
        result, err := handler(sctx, params)
        if err != nil {
            sctx.Errorf("Handler error: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Если результат не nil, сериализуем его в JSON
        if result != nil {
            w.Header().Set("Content-Type", "application/json")
            if err := json.NewEncoder(w).Encode(result); err != nil {
                sctx.Errorf("Error encoding response: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }
    }
}
