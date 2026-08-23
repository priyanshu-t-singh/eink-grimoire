package handlers

import (
	"le-grimoire/internal/models"
	"net/http"
)

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.RespondWithData(w, models.APIResponse[string]{Data: "OK"})
}
