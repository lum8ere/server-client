package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func (s *Service) uploadHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Println("Получен запрос на загрузку данных")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Printf("Ошибка чтения данных: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	latestFrame := fmt.Sprintf("%s/latest_frame.jpg", s.uploadPath)
	if err := os.WriteFile(latestFrame, data, 0644); err != nil {
		s.logger.Printf("Ошибка обновления latest_frame.jpg: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Service) uploadScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Println("Получен запрос на загрузку скриншота")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Printf("Ошибка чтения данных: %v", err)
		http.Error(w, "Ошибка чтения данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	filename := fmt.Sprintf("%s/screenshot_%d.jpg", s.uploadPath, time.Now().UnixNano())
	if err := os.WriteFile(filename, data, 0644); err != nil {
		s.logger.Printf("Ошибка сохранения файла: %v", err)
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}
	latestScreenshot := fmt.Sprintf("%s/latest_screenshot.jpg", s.uploadPath)
	if err := os.WriteFile(latestScreenshot, data, 0644); err != nil {
		s.logger.Printf("Ошибка обновления latest_screenshot.jpg: %v", err)
	}
	s.logger.Printf("Скриншот успешно сохранён: %s", filename)
	w.WriteHeader(http.StatusOK)
}
