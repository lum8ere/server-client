package rest_middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"

	"backed-api-v2/libs/5_common/smart_context"
)

// SmartHandlerFunc определяет обработчик, принимающий smart_context.
type SmartHandlerFunc func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request)

// WithRequestId извлекает из заголовка X-Request-Id или генерирует новый, и добавляет его в smart_context.
func WithRequestId(handler SmartHandlerFunc) SmartHandlerFunc {
	return func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-Request-Id")
		if requestId == "" {
			requestId = uuid.New().String()
		}
		sctx = sctx.WithField("request_id", requestId)
		w.Header().Set("X-Request-Id", requestId)
		handler(sctx, w, r)
	}
}

// WithSessionId извлекает из заголовка X-Session-Id или генерирует новый, и добавляет его в smart_context.
func WithSessionId(handler SmartHandlerFunc) SmartHandlerFunc {
	return func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get("X-Session-Id")
		if sessionId == "" {
			sessionId = uuid.New().String()
		}
		sctx = sctx.WithSessionId(sessionId)
		w.Header().Set("X-Session-Id", sessionId)
		handler(sctx, w, r)
	}
}

// WithDeviceId извлекает из заголовка X-Device-Identifier или генерирует новый, и добавляет его в smart_context.
func WithDeviceIdentifier(handler SmartHandlerFunc) SmartHandlerFunc {
	return func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
		deviceId := r.Header.Get("X-Device-Identifier") // МЫ СЧИТАЕМ ЧТО ОН ВСЕГД БУДЕТ
		if deviceId == "" {
			http.Error(w, "Missing X-Device-Identifier header", http.StatusBadRequest)
			return
		}
		sctx = sctx.WithDeviceIdentifier(deviceId)
		w.Header().Set("X-Device-Identifier", deviceId)
		handler(sctx, w, r)
	}
}

// WithRecoverer оборачивает обработчик, отлавливая возможные panics.
func WithRecoverer(handler SmartHandlerFunc) SmartHandlerFunc {
	return func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				sctx.Errorf("Recovered from panic: %v\nStack: %s", rec, debug.Stack())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		handler(sctx, w, r)
	}
}

// WithRestApiSmartContext объединяет цепочку middleware и возвращает стандартный http.HandlerFunc.
func WithRestApiSmartContext(sctx smart_context.ISmartContext, handler SmartHandlerFunc) http.HandlerFunc {
	// Собираем цепочку: WithRecoverer, затем WithSessionId, WithRequestId и WithDeviceIdentifier.
	chain := WithRecoverer(WithSessionId(WithRequestId(WithDeviceIdentifier(handler))))
	return func(w http.ResponseWriter, r *http.Request) {
		chain(sctx, w, r)
	}
}

// Аналогичная обёртка для WebSocket, если требуется.
func WithWsApiSmartContext(sctx smart_context.ISmartContext, handler SmartHandlerFunc) http.HandlerFunc {
	chain := WithRecoverer(WithSessionId(WithRequestId(WithDeviceIdentifier(handler))))
	return func(w http.ResponseWriter, r *http.Request) {
		chain(sctx, w, r)
	}
}

// WithWaitGroup добавляет обработку waitgroup.
func WithWaitGroup(smartHandler SmartHandlerFunc) SmartHandlerFunc {
	return func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
		// Проверяем, не закрыт ли контекст сервера.
		if err := sctx.GetContext().Err(); err != nil {
			sctx.Warnf("Server context is closed: %v. Cannot run request", err)
			http.Error(w, "Server context is closed", http.StatusServiceUnavailable)
			return
		}

		// Добавляем счетчик waitgroup, если он есть.
		if wg := sctx.GetWaitGroup(); wg != nil {
			wg.Add(1)
			defer wg.Done()
		}

		smartHandler(sctx, w, r)
	}
}
