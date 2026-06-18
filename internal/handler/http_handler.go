package handler

import (
	"encoding/json"
	"net/http"
)

type HTTPHandler struct {
	viewDataHandler ViewDataHandler
}

func NewHTTPHandler(h ViewDataHandler) *HTTPHandler {
	return &HTTPHandler{viewDataHandler: h}
}

func (h *HTTPHandler) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/view-data", h.handleViewData)

	return http.ListenAndServe(addr, mux)
}

func (h *HTTPHandler) handleViewData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response, err := h.viewDataHandler.GetViewData()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
