package rest_middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// RoleMiddleware проверяет, что токен (из заголовка Authorization) содержит роль,
// достаточную для доступа к данному ресурсу.
// Пример: для admin необходимо, чтобы роль была "admin".
func RoleMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(tokenStr, " ")
		if len(parts) != 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenStr = parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secret := []byte(getJWTSecret())
			return secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		userRole, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Role not found in token", http.StatusUnauthorized)
			return
		}
		if !roleSufficient(userRole, requiredRole) {
			http.Error(w, "Insufficient privileges", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func roleSufficient(userRole, requiredRole string) bool {
	rolePriority := map[string]int{
		"observer":  1,
		"observer+": 2,
		"admin":     3,
	}
	return rolePriority[userRole] >= rolePriority[requiredRole]
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "defaultsecret"
	}
	return secret
}
