package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Location struct {
	Status string  `json:"status"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
}

func (s *Service) mapHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")

	s.m.RLock()
	client, ok := s.clients[clientID]
	if !ok || client.IP == "" {
		http.Error(w, "Нет клиента с сохранённым публичным IP", http.StatusNotFound)
		return
	}
	s.m.RUnlock()

	geoURL := fmt.Sprintf("http://ip-api.com/json/%s", client.IP)
	resp, err := http.Get(geoURL)
	if err != nil {
		http.Error(w, "Ошибка запроса геоданных", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

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

func (s *Service) mapFrontendHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    clientID := vars["id"]

	s.m.RLock()
	client, ok := s.clients[clientID]
	if !ok || client.IP == "" {
		http.Error(w, "Нет клиента с сохранённым публичным IP", http.StatusNotFound)
		return
	}
	s.m.RUnlock()

	geoURL := fmt.Sprintf("http://ip-api.com/json/%s", client.IP)
	resp, err := http.Get(geoURL)
	if err != nil {
		http.Error(w, "Ошибка запроса геоданных", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var loc Location
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		http.Error(w, "Ошибка декодирования геоданных", http.StatusInternalServerError)
		return
	}
	if loc.Status != "success" {
		http.Error(w, "Геокодирование не удалось", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	locData, err := json.Marshal(loc)
	if err != nil {
		http.Error(w, "Ошибка сериализации", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
    w.Write(locData)
}