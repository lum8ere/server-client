package service

import (
	"encoding/json"
	"html/template"
	"maps"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID      string
	Conn    *websocket.Conn
	IP      string
	Status  string
	Metrics Metrics
	AppsServices string // Здесь будем сохранять JSON-строку с данными
}

func (s *Service) clientHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Не указан клиент", http.StatusBadRequest)
		return
	}
	tmpl, err := template.ParseFiles("static/client.html")
	if err != nil {
		s.logger.Printf("Ошибка загрузки шаблона client.html: %v", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, clientID)
}

func (s *Service) clientsHandler(w http.ResponseWriter, r *http.Request) {
	s.m.Lock()
	clients := slices.Collect(maps.Values(s.clients))
	s.m.Unlock()

	tmpl, err := template.ParseFiles("static/clients.html")
	if err != nil {
		s.logger.Printf("Ошибка загрузки шаблона clients.html: %v", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, clients)
}

func (s *Service) GetClients(w http.ResponseWriter, r *http.Request) {
	s.m.Lock()
	clients := slices.Collect(maps.Values(s.clients))
	s.m.Unlock()

	clientsData, err := json.Marshal(clients)
	if err != nil {
		http.Error(w, "Ошибка сериализации", http.StatusInternalServerError)
		return
	}
	w.Write(clientsData)
}

func (s *Service) clientAppsDataHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Не указан клиент", http.StatusBadRequest)
		return
	}

	s.m.RLock()
	client, ok := s.clients[clientID]
	s.m.RUnlock()
	if !ok {
		http.Error(w, "Клиент не найден", http.StatusNotFound)
		return
	}

	// Если клиент ещё не отправлял данные, возвращаем пустой JSON.
	if client.AppsServices == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"processes": [], "services": []}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(client.AppsServices))
}
type Metrics struct {
	DiskTotal       uint64 `json:"disk_total"`
	DiskFree        uint64 `json:"disk_free"`
	MemoryTotal     uint64 `json:"memory_total"`
	MemoryAvailable uint64 `json:"memory_available"`
	Processor       string `json:"processor"`
	OS              string `json:"os"`
	HasPassword     bool   `json:"has_password"`
	MinimumPasswordLenght int `json:"minimum_password_lenght"`
	Error           string `json:"error,omitempty"`
	PcName 			string `json:"pc_name,omitempty"`
}

var nilMetrics Metrics

func (s *Service) clientMetricsHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Не указан клиент", http.StatusBadRequest)
		return
	}

	client, ok := s.clients[clientID]
	if !ok {
		http.Error(w, "Клиент не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if client.Metrics == nilMetrics {
		w.Write([]byte("{}"))
		return
	}
	jsonData, err := json.Marshal(client.Metrics)
	if err != nil {
		http.Error(w, "Ошибка сериализации", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
