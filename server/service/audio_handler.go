package service

import "github.com/gorilla/websocket"

func (s *Service) HandleAudioSource(clientID string, conn *websocket.Conn) {
	defer func() {
		s.m.Lock()
		delete(s.audioSources, clientID)
		s.m.Unlock()
		conn.Close()
		s.logger.Printf("Audio source disconnected: %s", clientID)
	}()

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			s.logger.Printf("Audio source read error (%s): %v", clientID, err)
			return
		}
		if msgType == websocket.BinaryMessage {
			// Это PCM-данные
			s.broadcastAudio(clientID, data)
		} else {
			s.logger.Printf("Audio source (%s) non-binary msg received", clientID)
		}
	}
}

func (s *Service) broadcastAudio(clientID string, pcm []byte) {
    s.m.RLock()
    listeners := s.audioListeners[clientID]
    s.m.RUnlock()

    for _, lconn := range listeners {
        err := lconn.WriteMessage(websocket.BinaryMessage, pcm)
        if err != nil {
            s.logger.Printf("broadcastAudio error -> removing listener: %v", err)
            // можно закрыть сокет и убрать из списка
        }
    }
}
