package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/service"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	uploadPath = "./uploads"
)







func main() {
	logger := log.New(os.Stdout, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		logger.Fatalf("Ошибка создания директории: %v", err)
	}

	go startJPGCleanupJob(10 * time.Minute, logger)

	s := service.New(logger, uploadPath)

	r := mux.NewRouter()

	service.RegisterRoutes(r, s)

    srv := &http.Server{
        Handler:      r,
        Addr:         "127.0.0.1:4000",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

	logger.Println("Сервер запущен на :4000")

	log.Fatal(srv.ListenAndServe())
}



func startJPGCleanupJob(interval time.Duration,logger *log.Logger) {
	for {
		time.Sleep(interval)
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
}


