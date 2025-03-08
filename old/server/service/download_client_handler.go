package service

import (
	"net/http"
	"path/filepath"
)

func (s *Service) downloadClientHandler(w http.ResponseWriter, r *http.Request) {
	// Формируем полный путь до файла client.exe в папке uploads
	clientPath := filepath.Join(s.uploadPath, "client.exe")

	// Устанавливаем заголовки для скачивания файла как attachment
	w.Header().Set("Content-Disposition", "attachment; filename=client.exe")
	w.Header().Set("Content-Type", "application/octet-stream")

	// Отдаем файл
	http.ServeFile(w, r, clientPath)
}
