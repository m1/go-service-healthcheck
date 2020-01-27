package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/m1/go-service-healthcheck/response"
)

// InternalHandler ...
type InternalHandler struct {
	*API
}

// NewDebugHandler ...
func NewInternalHandler(api *API) *InternalHandler {
	return &InternalHandler{api}
}

// GetRoutes ...
func (h *InternalHandler) GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/health", h.GetHealth)
	return router
}

func (h *InternalHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	response.RespondOk(w, r)
}
