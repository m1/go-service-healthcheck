package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"

	"github.com/m1/go-service-healthcheck/response"
)

const ParamServiceID = "serviceID"

// InternalHandler ...
type ServicesHandler struct {
	*API
}

// NewDebugHandler ...
func NewServicesHandler(api *API) *ServicesHandler {
	return &ServicesHandler{api}
}

// GetRoutes ...
func (h *ServicesHandler) GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", h.GetServices)
	router.Get(fmt.Sprintf("/{%v}/events", ParamServiceID), h.GetEvents)
	return router
}

func (h *ServicesHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.Repositories.Service.GetAll()
	if err != nil {
		response.RespondInternalServerError(w, r, err)
		return
	}
	response.RespondOk(w, r, &response.Data{Data: &services})
}

func (h *ServicesHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	serviceID, err := strconv.Atoi(chi.URLParam(r, ParamServiceID))
	if err != nil {
		response.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	service, err := h.Repositories.Service.Get(serviceID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.RespondError(w, r, http.StatusNotFound, err)
			return
		}
		response.RespondInternalServerError(w, r, err)
		return
	}

	events, err := h.Repositories.Service.GetEvents(service.ID)
	if err != nil {
		response.RespondInternalServerError(w, r, err)
		return
	}

	response.RespondOk(w, r, &response.Data{Data: &events})
}
