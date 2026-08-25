package handlers

import (
	"fmt"
	"le-grimoire/internal/middleware"
	"net/http"
)

func (h *Handler) CurrentPageHandler(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := middleware.GetDeviceID(r.Context())
	if !ok {
		h.RespondWithStatusError(w, http.StatusUnauthorized, fmt.Errorf("device ID not found in headers"))
		return
	}

	ds, err := h.App.DeviceRepository.GetDeviceState(deviceID) // TODO: confirm real accessor path
	if err != nil {
		h.RespondWithError(w, err)
		return
	}
	h.RespondWithData(w, ds)
}
