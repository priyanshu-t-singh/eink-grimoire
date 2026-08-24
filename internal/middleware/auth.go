package middleware

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

	"le-grimoire/internal/device"
)

type contextKey string

const DeviceIDKey contextKey = "device_id"

func DeviceAuth(repo *device.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			deviceID := r.Header.Get("X-Device-Id")
			authHeader := r.Header.Get("Authorization")

			if deviceID == "" || authHeader == "" {
				http.Error(w, "Missing authentication headers", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			apiKey := parts[1]

			// Use the repo directly instead of app.DeviceRepo
			expectedHash, err := repo.GetDeviceAuthHash(deviceID)
			if err != nil {
				log.Printf("Auth failed: device %s not found or DB error: %v", deviceID, err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			hash := sha256.Sum256([]byte(apiKey))
			incomingHash := hex.EncodeToString(hash[:])

			if subtle.ConstantTimeCompare([]byte(incomingHash), []byte(expectedHash)) != 1 {
				log.Printf("Auth failed: invalid API key for device %s", deviceID)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), DeviceIDKey, deviceID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDeviceID(ctx context.Context) string {
	if id, ok := ctx.Value(DeviceIDKey).(string); ok {
		return id
	}
	return ""
}
