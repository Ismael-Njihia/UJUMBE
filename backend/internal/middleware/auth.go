package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/Ismael-Njihia/UJUMBE/backend/internal/database"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func AuthMiddleware(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				// Try Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
					apiKey = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			if apiKey == "" {
				http.Error(w, `{"error":"API key required"}`, http.StatusUnauthorized)
				return
			}

			var userID uuid.UUID
			err := db.QueryRow("SELECT id FROM users WHERE api_key = $1", apiKey).Scan(&userID)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, `{"error":"Invalid API key"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	return userID, ok
}
