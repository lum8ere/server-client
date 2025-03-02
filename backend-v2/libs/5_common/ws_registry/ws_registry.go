package ws_registry

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients      = make(map[string]*websocket.Conn) // ключ – deviceId
	clientsMutex sync.RWMutex
)

// SetClient сохраняет соединение для указанного deviceId.
func SetClient(deviceId string, conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients[deviceId] = conn
}

// RemoveClient удаляет соединение по deviceId.
func RemoveClient(deviceId string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	delete(clients, deviceId)
}

// GetClient возвращает соединение для указанного deviceId.
func GetClient(deviceId string) (*websocket.Conn, bool) {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()
	conn, ok := clients[deviceId]
	return conn, ok
}
