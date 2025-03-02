package ws_server

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/env_vars"
	"backed-api-v2/libs/5_common/safe_go"
	"backed-api-v2/libs/5_common/smart_context"
	"compress/flate"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
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

func (u *WsUpgrader) sendResponse(sctx smart_context.ISmartContext, messageChan chan<- Message, msg Message) {
	// если контекст прервали то не нужно отправять ответ!
	msg.Sctx = sctx // чисто для логгирования!
	ctx := sctx.GetContext()

	// Use a select statement to handle context cancellation or sending the message.
	select {
	case <-ctx.Done():
		// If the context is canceled, log the cancellation and return without sending the message.
		sctx.Warnf("WebSocket connection: sendResponse: context canceled, when was trying to send '%s' to out channel", messageTypeToString(msg.Type))
		return
	case messageChan <- msg:
		// If the context is active, send the response to the channel.
		// sctx.Debug("sendResponse(): response sent successfully to out channel")
	}
}

func (u *WsUpgrader) HandleWebSocket(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	serverCtx := sctx.GetContext()

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

	sessionCtx, cancelSessionCtx := context.WithCancel(r.Context())
	defer cancelSessionCtx()
	sctx = sctx.WithContext(sessionCtx)

	var regErr error
	sctx, regErr = registerDevice(sctx)
	if regErr != nil {
		sctx.Errorf("Device registration failed: %v", regErr)
		conn.Close()
		return
	}

	defer func() {
		sctx.Infof("WebSocket connection: defer block: closing messages channel")
		close(messageChan)
	}()

	safe_go.SafeGo(sctx, func() {
		// Канал запроектся только в defer когда все писатели запишут туда что хотели и завершат свою работу (wg опустет)
		u.writePump(conn, messageChan, sctx)
	})

	isRealEnv := smart_context.IsRealEnv()
	if isRealEnv {
		// Set initial read deadline
		err = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			sctx.Errorf("WebSocket connection: error setting initial read deadline: %v", err)
		}
	}

	// Set pong handler to extend the deadline
	conn.SetPongHandler(func(_ string) error {
		sctx.Debugf("WebSocket connection: received %s", messageTypeToString(websocket.PongMessage))
		if isRealEnv {
			err = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			if err != nil {
				sctx.Errorf("WebSocket connection: error extending read deadline: %v", err)
			}
		}
		return nil
	})

	// Set ping handler to log pings
	conn.SetPingHandler(func(appData string) error {
		// in case client sends a ping message, we should respond with a pong message
		sctx.Debugf("WebSocket connection: received %s", messageTypeToString(websocket.PingMessage))
		msg := Message{
			Type: websocket.PongMessage,
			Data: []byte(appData),
		}
		u.sendResponse(sctx, messageChan, msg)
		// sctx.Debugf("WebSocket connection: pong sent")
		return nil
	})

	// Set close handler to log close messages
	conn.SetCloseHandler(func(code int, text string) error {
		sctx.Infof("WebSocket connection: connection closed (%d - %s)", code, text)
		// Если соединение закрыто клиентом, обновляем статус устройства на OFFLINE.
		if code == websocket.CloseNormalClosure || code == websocket.CloseGoingAway {
			setDeviceStatusOffline(sctx)
		}
		return nil
	})

	// Ping the client periodically
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Ping every 30 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				msg := Message{
					Type: websocket.PingMessage,
					Data: nil,
				}
				u.sendResponse(sctx, messageChan, msg)
				// sctx.Debugf("WebSocket connection: ping sent")

			case <-sessionCtx.Done():
				// юзер закрыл браузер (или сервер закрылся и мы дождались завершения всех запросов и отправки их ответов) - перестаем слать пинги
				return
			}
		}
	}()

	go func() {
		for {
			messageType, messageData, err := conn.ReadMessage()
			if err != nil {
				// Check if the error is a result of a normal closure
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					// тут ошибка потому что не должен фронтед так рвать сообщение
					// UPDATE: но если я браузер закрою то это нормально - поэтому понижаем до Warn
					sctx.Warnf("WebSocket connection: read message loop: closed abnormaly by client: %v", err.Error())
				} else {
					// This could be a normal disconnection or a different close code that does not signify an error
					sctx.Infof("WebSocket connection: read message loop: closed normally by client: %v", err.Error())
				}

				cancelSessionCtx() // закрываем контекст - чтобы новые сообщения прекратить слать
				sctx.Infof("WebSocket connection: read message loop: cancelled session context, exiting read message loop")
				break
			}

			sctx.Debugf("Received %s", messageTypeToString(messageType))
			if messageType == websocket.TextMessage {
				// Обрабатываем текстовое сообщение: если начинается с "register:", то вызываем processWsMessage.
				processWsMessage(sctx, string(messageData))
			} else {
				// Для других типов сообщений можно добавить дополнительную обработку
				sctx.Infof("Non-text message received")
			}

			sctx.Debugf("WebSocket connection: read message loop: received %s", messageTypeToString(messageType))

			if serverCtx.Err() != nil {
				// Если контекст прервали - то не нужно обрабатывать запросы - просто их игнорируем
				sctx.Infof("WebSocket connection: read message loop: server context canceled, need to respond with error")
				// done: нельзя просто игнорировать- надо вернуть ошибку!
				// u.safeProcessRequest(sctx, messageChan, msg, true)
				continue // продолжаем дальше вычитывать сообщения - коннект скоро закроют и мы выйдем из цикла
			}

			wg.Add(1)
			safe_go.SafeGo(sctx, func() {
				defer wg.Done()
				// если контекст прервут мы быстро выйдем даже если результат уже готов для отправки в канал и наружу
				// TODO: не надо выходить по прерыванию контекста
				// u.safeProcessRequest(sctx, messageChan, msg, false)
			})
		}
	}()

	select {
	case <-serverCtx.Done():
		// если закроют serverCtx то ждем завершения всех запросов и прерываем сессию
		sctx.Infof("WebSocket connection: main block: server context canceled, waiting all requests to finish")

		// TODO: Все-таки нет смысла держать соединения открытыми если сервер закрывается
		// в эти соединения клиенты все равно уже не смогут отправить новые запросы - мы их отклоняем
		// а в таком случае нет смысла держать соединение - клиет конечно может получить ответ но нужен ли он ему?
		// может быть ему лучше переконнектиться и получить новый сокет?
		// ИТОГО: Надо рвать связь сразу т.к. фронтед всё равно может отправить новые запросы
		wg.Wait() // ждем завершения всех запросов - новые уже не возникают
		sctx.Infof("WebSocket connection: main block: all requests finished - send CloseMessage to user")

		// send CloseMessage to user
		msg := Message{
			Type: websocket.CloseMessage,
			Data: websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server context canceled"),
		}
		u.sendResponse(sctx, messageChan, msg)
		// после этого сообщения юзер должен разорвать соединение. но если не разорвет - не страшно - мы все равно уже выходим

		sctx.Infof("WebSocket connection: main block: CloseMessage sent to user, canceling session context")
		cancelSessionCtx() // закрываем контекст - чтобы новые сообщения прекратить слать - закрываем только здесь тк выше слали CloseMessage
		sctx.Infof("WebSocket connection: main block: session context canceled, exiting main block")

	case <-sessionCtx.Done():
		// сюда попадаем только после закрытия соединения юзером
		// уже перестали слать ответы (sessionCtx это предотвращает) но дождаться завершения запросов нужно
		sctx.Infof("WebSocket connection: main block: session context canceled, waiting all requests to finish")
		wg.Wait() // ждем завершения всех запросов - новые уже не возникают
		sctx.Infof("WebSocket connection: main block: all requests finished, exiting main block")
	}
}

func messageTypeToString(messageType int) string {
	switch messageType {
	case websocket.TextMessage:
		return "TextMessage"
	case websocket.BinaryMessage:
		return "BinaryMessage"
	case websocket.CloseMessage:
		return "CloseMessage"
	case websocket.PingMessage:
		return "PingMessage"
	case websocket.PongMessage:
		return "PongMessage"
	default:
		return "Unknown"
	}
}

func (u *WsUpgrader) writePump(conn *websocket.Conn, messageChan chan Message, sctx smart_context.ISmartContext) {
	defer func() {
		sctx.Infof("WebSocket connection: close connection on server side")
		conn.Close() // закрываем соединение когда выйдем из функции
	}()

	// если контекст прервут то все быстро прекратят слать данные в канал, все запросы завершатся и мы выйдем и канал закроем
	// поэтому тут безопасно и красиво всё вычитывать до конца
	for msg := range messageChan {
		startSendingAt := time.Now()
		curSctx := sctx
		if msg.Sctx != nil {
			curSctx = msg.Sctx
		}

		messageType := messageTypeToString(msg.Type)
		var err error

		switch msg.Type {
		case websocket.PingMessage, websocket.PongMessage, websocket.CloseMessage:
			// Send a control message using WriteControl
			deadline := time.Second
			if !smart_context.IsRealEnv() {
				deadline = time.Hour // в локальном режиме дедлайн час , чтобы при отладке не рвало соединение
			}

			err = conn.WriteControl(msg.Type, msg.Data, time.Now().Add(deadline))

		default:

			// Send a regular message using WriteMessage
			err = conn.WriteMessage(msg.Type, msg.Data)
		}

		duration := time.Since(startSendingAt)
		size := len(msg.Data)

		if err != nil {
			// эти статусы ошибки не ошибка сервера (например, клиент закрыл соединение), для них функция вернет false
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				curSctx.
					LogField("size", size).
					LogField("duration", duration).
					Errorf("WebSocket connection: error sending %s: %v", messageType, err)
			} else {
				curSctx.
					LogField("size", size).
					LogField("duration", duration).
					Debugf("WebSocket connection: skip sending %s: %v", messageType, err)
			}
		} else {
			curSctx.
				LogField("size", size).
				LogField("duration", duration).
				Debugf("WebSocket connection: sent %s", messageType)
		}
	}
}

// processWsMessage обрабатывает входящие текстовые сообщения.
func processWsMessage(sctx smart_context.ISmartContext, msg string) {
	// Извлекаем фактический device_id, сохранённый при регистрации
	deviceDBID, ok := sctx.GetDataFields().GetField("device_id")
	if !ok || deviceDBID == "" {
		sctx.Warn("No device_id found in context, cannot process message")
		return
	}
	// Приводим к строке
	realDeviceID, _ := deviceDBID.(string)

	// Попытка распарсить сообщение как метрики.
	var m model.Metric
	if err := json.Unmarshal([]byte(msg), &m); err == nil {
		if m.DeviceID == "" {
			m.DeviceID = realDeviceID
		}
		m.CreatedAt = time.Now()
		if err := sctx.GetDB().Create(&m).Error; err != nil {
			sctx.Errorf("Error saving metrics for device %s: %v", m.DeviceID, err)
			return
		}
		sctx.Infof("Metrics saved for device: %s", m.DeviceID)
		return
	}
	// Если сообщение не метрики, просто логируем его.
	sctx.Infof("Received WS message: %s", msg)
}

func registerDevice(sctx smart_context.ISmartContext) (smart_context.ISmartContext, error) {
	deviceIdentifier := sctx.GetDeviceIdentifier() // device_identifier из заголовка
	if deviceIdentifier == "" {
		sctx.Warnf("No device identifier provided in context")
		return sctx, fmt.Errorf("missing device identifier")
	}

	var device model.Device
	db := sctx.GetDB()
	err := db.Where("device_identifier = ?", deviceIdentifier).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Устройство не найдено, создаём новую запись
			device = model.Device{
				DeviceIdentifier: deviceIdentifier,
				Status:           "ONLINE",
				UserID:           "411c0868-638a-4370-ad42-650cab9fb227", // временно
				LastSeen:         time.Now(),
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			if err := db.Create(&device).Error; err != nil {
				sctx.Errorf("Error registering device %s: %v", deviceIdentifier, err)
				return sctx, err
			}
			sctx.Infof("Device registered: %s", deviceIdentifier)
		} else {
			sctx.Errorf("DB error when processing device %s: %v", deviceIdentifier, err)
			return sctx, err
		}
	} else {
		// Устройство найдено, обновляем информацию
		device.LastSeen = time.Now()
		device.Status = "ONLINE"
		if err := db.Save(&device).Error; err != nil {
			sctx.Errorf("Error updating device %s: %v", deviceIdentifier, err)
			return sctx, err
		}
		sctx.Infof("Device updated: %s", deviceIdentifier)
	}
	// Сохраняем фактический device_id (primary key) в контекст под новым ключом
	sctx = sctx.WithField("device_id", device.ID)
	sctx.Infof("Saved device_id in context: %v", device.ID)
	return sctx, nil
}

func setDeviceStatusOffline(sctx smart_context.ISmartContext) error {
	deviceIdentifier := sctx.GetDeviceIdentifier() // device_identifier из заголовка
	if deviceIdentifier == "" {
		sctx.Warnf("No device identifier provided in context")
		return nil
	}
	if err := sctx.GetDB().Model(&model.Device{}).Where("device_identifier = ?", deviceIdentifier).Update("status", "OFFLINE").Error; err != nil {
		sctx.Errorf("Error updating device %s status to OFFLINE: %v", deviceIdentifier, err)
	} else {
		sctx.Infof("Device %s set to OFFLINE", deviceIdentifier)
	}

	return nil
}
