package rest_middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const SessionIDKey = "session_id"

func GenerateSessionID() string {
	// Generate a new UUID
	id := uuid.New()
	// Convert UUID to string and return as SessionID
	return id.String()
}

// Middleware to create and store the logger in the context
func SessionID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session ID from the header
		sessionID := r.Header.Get("X-Session-Id")
		if sessionID == "" {
			sessionID = GenerateSessionID()
		}

		// Store the session ID in the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, SessionIDKey, sessionID)

		// Call the next handler, passing the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSessionID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if sessionID, ok := ctx.Value(SessionIDKey).(string); ok {
		return sessionID
	}
	return ""
}
