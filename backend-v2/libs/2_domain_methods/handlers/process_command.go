package handlers

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"backed-api-v2/libs/5_common/ws_registry"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func SendCommandHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
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

	// Сохраним команду в БД со статусом "pending"
	cmdRecord := &model.Command{
		DeviceID:    deviceId,
		CommandType: command,
		UserID:      "a7c4265d-545c-404f-a1ef-daf4af2dfb12", // или другой способ идентификации инициатора TODO: После подготовки пользователей, сюда надо будет вставить его
		Status:      "PENDING",
		CreatedAt:   time.Now(),
	}

	if err := sctx.GetDB().Create(cmdRecord).Error; err != nil {
		return nil, fmt.Errorf("error saving command to db: %w", err)
	}
	sctx.Infof("Command saved with ID: %s", cmdRecord.ID)

	// Находим соединение для этого устройства из глобального реестра WebSocket-соединений.
	conn, ok := ws_registry.GetClient(deviceId)
	if !ok {
		sctx.Infof("Client with device '%s' not found, command stored for later execution", deviceId)
		return map[string]any{"status": "PENDING"}, nil
	}

	// Отправляем команду по WebSocket
	if err := conn.WriteMessage(websocket.TextMessage, []byte(command)); err != nil {
		sctx.GetDB().Model(cmdRecord).Update("status", "ERROR")
		return nil, fmt.Errorf("error sending command: %w", err)
	}

	sctx.Infof("Command '%s' sent to device '%s'", command, deviceId)
	sctx.GetDB().Model(cmdRecord).Updates(map[string]any{
		"status":      "SENT",
		"executed_at": time.Now(),
	})

	resp := types.ANY_DATA{
		"status":  "success",
		"device":  deviceId,
		"command": command,
	}

	return resp, nil
}
