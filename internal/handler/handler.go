package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PaulusBilly/internship-test-API/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) GetViewData(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.BuildResponse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
