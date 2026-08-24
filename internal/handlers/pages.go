package handlers

import (
	"le-grimoire/internal/middleware"
	"net/http"
)

func (h *Handler) CurrentPageHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := middleware.GetDeviceID(r.Context())
	if deviceID == "" {
		http.Error(w, "Device ID not found in context", http.StatusUnauthorized)
		return
	}

	// TODO: Implement the logic to retrieve the current page information
	h.RespondWithData(w, "Current page information")
}
