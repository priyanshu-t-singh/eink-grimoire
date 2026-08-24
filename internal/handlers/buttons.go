package handlers

import (
	"le-grimoire/internal/middleware"
	"net/http"
)

type buttonRequest struct {
	ButtonID string `json:"button_id"` // Expected values: "A", "B", "C", "D", "E", "F"
	Type     string `json:"type"`      // Expected values: "short_press", "long_press"
}

func (h *Handler) PushButtonHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := middleware.GetDeviceID(r.Context())
	if deviceID == "" {
		http.Error(w, "Device ID not found in context", http.StatusUnauthorized)
		return
	}

	// TODO: Implement the logic to handle the push button event
	h.RespondWithData(w, "Button pushed successfully")
}
