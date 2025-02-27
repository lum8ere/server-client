package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var 	upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ClientMessage struct {
	Command string          `json:"command"`
	Data    json.RawMessage `json:"data"`
}

func (s *Service)wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Println("Ошибка апгрейда WebSocket:", err)
		return
	}

	s.clientCounter.Add(1)
	clientID := fmt.Sprintf("client-%d", s.clientCounter.Load())
	client := Client{ID: clientID, Conn: conn, IP: ""}
	
	s.m.Lock()
	s.clients[clientID] = client
	s.m.Unlock()

	s.logger.Printf("WebSocket клиент %s подключён", clientID)

	if err := conn.WriteMessage(websocket.TextMessage, []byte("id:"+clientID)); err != nil {
		s.logger.Printf("Ошибка отправки ID клиенту %s: %v", clientID, err)
	}

	defer func() {
		s.m.Lock()
		conn.Close()
		delete(s.clients, clientID)
		s.m.Unlock()
		s.logger.Printf("WebSocket клиент %s отключён", clientID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.logger.Printf("Ошибка чтения WebSocket у клиента %s: %v", clientID, err)
			break
		}
		msgStr := string(message)
		s.logger.Printf("Клиент %s отправил: %s", clientID, msgStr)
		if strings.HasPrefix(msgStr, "ip:") {
			ip := strings.TrimPrefix(msgStr, "ip:")
			client.IP = ip
			s.logger.Printf("Сохранён публичный IP для клиента %s: %s", clientID, ip)

		} else if strings.HasPrefix(msgStr, "{") {
			var cm ClientMessage
			if err := json.Unmarshal(message, &cm); err != nil {
				s.logger.Printf("Ошибка декодирования JSON от клиента %s: %v", clientID, err)
				continue
			}

			switch cm.Command {
			case "metrics":
				var m Metrics
				if err := json.Unmarshal(cm.Data, &m); err != nil {
					s.logger.Printf("Ошибка декодирования метрик от клиента %s: %v", clientID, err)
				}
				client.Metrics = m
				s.logger.Printf("Получены метрики от клиента %s: %+v", clientID, m)

			default:
				s.logger.Printf("Неизвестная команда JSON от клиента %s: %s", clientID, cm.Command)
			}
		} else {
			s.logger.Printf("Получено сообщение от клиента %s: %s", clientID, msgStr)
		}
	}
}

func(s *Service) commandHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")
	if cmd == "" {
		http.Error(w, "Параметр cmd отсутствует", http.StatusBadRequest)
		return
	}
	targetID := r.URL.Query().Get("id")

	s.m.RLock()
	defer s.m.RUnlock()
	if targetID != "" {
		client, ok := s.clients[targetID]
		if !ok {
			http.Error(w, "Клиент с таким ID не найден", http.StatusNotFound)
			return
		}
		if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(cmd)); err != nil {
			s.logger.Printf("Ошибка отправки команды клиенту %s: %v", targetID, err)
			http.Error(w, "Ошибка отправки команды", http.StatusInternalServerError)
			return
		}
		s.logger.Printf("Команда '%s' отправлена клиенту %s", cmd, targetID)
	} else {

		for _, client := range s.clients {
			if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(cmd)); err != nil {
				s.logger.Printf("Ошибка отправки команды клиенту %s: %v", client.ID, err)
			} else {
				s.logger.Printf("Команда '%s' отправлена клиенту %s", cmd, client.ID)
			}
		}
	}
	w.Write([]byte("Команда отправлена"))
}