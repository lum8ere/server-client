package auth

import (
	"encoding/json"
	"net/http"

	"backed-api-v2/libs/5_common/smart_context"
)

// LoginRequest – структура запроса на логин.
type LoginRequest struct {
	ProviderCode string `json:"ProviderCode"`
	Login        string `json:"Login"`
	Password     string `json:"Password"`
	// Можно добавить Token, DeviceInfo и другие поля по необходимости.
}

// LoginResponse – структура ответа при успешном логине.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	// Дополнительные поля, если нужны.
}

// LoginHandler обрабатывает POST запросы на /auth/login.
func LoginHandler(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sctx.Warnf("Failed to decode login request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Вызываем функцию логина (реализованную ниже).
	resp, err := Login(sctx, req)
	if err != nil {
		sctx.Warnf("Login failed: %v", err)
		// Не выдаем детальную информацию во фронтенд.
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		sctx.Warnf("Failed to encode login response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
