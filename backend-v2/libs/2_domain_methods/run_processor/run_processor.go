package run_processor

import (
	"backed-api-v2/libs/5_common/rest_middleware"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"encoding/json"
	"net/http"
	"strings"
)

// SmartHandlerFunc принимает smart-контекст и параметры запроса, возвращает объект (любой тип) и ошибку. Если объект не nil, он будет сериализован в JSON.
type SmartHandlerFunc func(sctx smart_context.ISmartContext, args types.ANY_DATA) (*types.ANY_DATA, error)

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

func WrapRestApiSmartHandler(sctx smart_context.ISmartContext, handler SmartHandlerFunc) http.HandlerFunc {
	// Сначала получаем базовый http.HandlerFunc, который извлекает параметры и вызывает handler.
	baseHandler := WrapSmartHandler(sctx, handler)
	// Затем оборачиваем его с помощью существующей цепочки middleware.
	// WithRestApiSmartContext ожидает SmartHandlerFunc с сигнатурой (sctx, w, r), поэтому создаём адаптер.
	return rest_middleware.WithRestApiSmartContext(sctx, func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
		baseHandler(w, r)
	})
}
