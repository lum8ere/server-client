package handlers

import (
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"backed-api-v2/libs/5_common/ws_registry"
	"fmt"

	"github.com/gorilla/websocket"
)

func SendCommandHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (*types.ANY_DATA, error) {
	sctx.Infof("SendCommandHandler started with args: %v", args)

	// Извлекаем device identifier из args
	deviceId, ok := args.GetStringValue("device_id")
	if !ok || deviceId == "" {
		return nil, fmt.Errorf("missing device identifier")
	}

	// Извлекаем команду, которую нужно отправить
	command, ok := args.GetStringValue("command")
	if !ok || command == "" {
		return nil, fmt.Errorf("missing command")
	}

	// Находим соединение для этого устройства из глобального реестра WebSocket-соединений.
	conn, found := ws_registry.GetClient(deviceId)
	if !found {
		return nil, fmt.Errorf("client with device '%s' not found", deviceId)
	}

	// Отправляем команду по WebSocket
	if err := conn.WriteMessage(websocket.TextMessage, []byte(command)); err != nil {
		return nil, fmt.Errorf("error sending command: %w", err)
	}
	sctx.Infof("Command '%s' sent to device '%s'", command, deviceId)

	// Формируем ответ, который будет сериализован в JSON
	resp := types.ANY_DATA{
		"status":  "success",
		"device":  deviceId,
		"command": command,
	}
	return types.AnyDataRef(resp), nil
}
