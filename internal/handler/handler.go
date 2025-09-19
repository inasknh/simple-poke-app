package handler

import (
	"encoding/json"
	"github.com/inasknh/simple-poke-app/internal/service"
	"net/http"
)

// Handler struct handles HTTP requests related to simple-poke-app.
type Handler struct {
	service service.Service
}

// NewHandler creates a new Handler instance with the given service.
func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SyncData(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	err := h.service.SyncData(ctx)
	if err != nil {
		httpResponseWrite(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	httpResponseWrite(rw, "OK", http.StatusOK)

}

func (h *Handler) GetItems(rw http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetItems(r.Context())
	if err != nil {
		httpResponseWrite(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	httpResponseWrite(rw, res, http.StatusOK)

}

// httpResponseWrite is a helper function to write JSON responses with the given data and status code.
func httpResponseWrite(rw http.ResponseWriter, data interface{}, statusCode int) {
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(statusCode)
	if data != nil {
		_ = json.NewEncoder(rw).Encode(data)
	}
}
