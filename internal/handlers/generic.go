package handlers

import (
	"net/http"
)

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.RespondWithData(w, "OK")
}
