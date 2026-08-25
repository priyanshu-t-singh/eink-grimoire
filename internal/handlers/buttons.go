package handlers

import (
	"encoding/json"
	"fmt"
	"le-grimoire/internal/middleware"
	"le-grimoire/internal/state"
	"net/http"
)

type buttonRequest struct {
	ButtonID string `json:"button_id"` // Expected values: "A", "B", "C", "D", "E", "F"
	Type     string `json:"type"`      // Expected values: "short_press", "long_press"
}

type pushButtonResponse struct {
	Action string            `json:"action"` // human-readable description, tracer-only
	State  state.DeviceState `json:"state"`
}

func (h *Handler) PushButtonHandler(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := middleware.GetDeviceID(r.Context())
	if !ok {
		h.RespondWithStatusError(w, http.StatusUnauthorized, fmt.Errorf("device ID not found in headers"))
		return
	}

	var req buttonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondWithStatusError(w, http.StatusBadRequest, err)
		return
	}

	ds, err := h.App.DeviceRepository.GetDeviceState(deviceID) // TODO: confirm real accessor path
	if err != nil {
		h.RespondWithError(w, err)
		return
	}

	action := state.ApplyButton(ds, req.ButtonID, req.Type)

	if err := h.App.DeviceRepository.SaveDeviceState(ds); err != nil {
		h.RespondWithError(w, err)
		return
	}

	h.RespondWithData(w, pushButtonResponse{Action: action, State: *ds})
}
