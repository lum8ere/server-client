package ws_server

import (
	"backed-api-v2/libs/4_common/env_vars"
	"backed-api-v2/libs/4_common/rest_middleware"
	"backed-api-v2/libs/4_common/smart_context"
	"compress/flate"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sctx smart_context.ISmartContext // чтобы логгировать под тем же request id
	Type int                         // WebSocket message type (e.g., TextMessage, PingMessage)
	Data []byte                      // Actual message data
}

type WsUpgrader struct {
	sctx     smart_context.ISmartContext
	Upgrader websocket.Upgrader

	apiRequestTimeout time.Duration
}

func NewWsUpgrader(sctx smart_context.ISmartContext) *WsUpgrader {
	apiRequestTimeoutSec := env_vars.GetEnvAsInt(sctx, "API_REQUEST_TIMEOUT", 600)

	return &WsUpgrader{
		sctx: sctx.
			LogField("component", "ws"),
		Upgrader: websocket.Upgrader{
			EnableCompression: true, // Enable Per-Message Deflate compression

			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		apiRequestTimeout: time.Duration(apiRequestTimeoutSec) * time.Second,
	}
}

func (u *WsUpgrader) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	// Log the request
	u.sctx.Infof("Unhandled request: %s %s", r.Method, r.URL.Path)

	// Return a 404 response
	http.NotFound(w, r)
}

func (u *WsUpgrader) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionId := rest_middleware.GetSessionID(ctx)

	sctx := u.sctx.
		WithSessionId(sessionId).         // для HandlerRun
		LogField("session_id", sessionId) // для логов

	conn, err := u.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		sctx.Errorf("WebSocket connection: request received, error upgrading to WebSocket: ", err)
		http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	// Optionally set compression level for outgoing messages
	if err := conn.SetCompressionLevel(flate.BestCompression); err != nil {
		sctx.Errorf("WebSocket connection: error setting compression level: %v", err)
	}

	// Session started
	u.sctx.Infof("WebSocket connection: request received, upgraded to WebSocket")

	messageChan := make(chan Message, 100) // Buffered channel to avoid blocking
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)
	sctx = sctx.WithContext(ctx)

	defer func() {
		sctx.Infof("WebSocket connection: cancel context, waiting all requests to finish")
		cancel()  // закрываем контекст - чтобы все запросы сразу прервались
		wg.Wait() // будем ждать пока ВСЕ запущенные запросы завершатся и их ответы прочитаются и отравятся юзеру - только потом закрываем канал
		sctx.Infof("WebSocket connection: all requests finished - close channel on server side")
		close(messageChan) // Close the message channel т.к. никто туда больше не пишет
	}()
}
