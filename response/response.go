package response

import (
	"net/http"

	"github.com/go-chi/render"
)

// Data is the parent model for returning data in this api,
// includes meta for pagination
type Data struct {
	Data       Model           `json:"data"`
	Pagination *PaginationData `json:"pagination,omitempty"`
}

// PaginationData ...
type PaginationData struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// Response ...
type Response struct {
	Status     *int             `json:"status"`
	StatusDesc *string          `json:"status_desc,"`
	Errors     *[]string        `json:"errors,omitempty"`
	Data       map[string]*Data `json:"data,omitempty"`
}

// Respond ...
func Respond(w http.ResponseWriter, r *http.Request, status int, datum ...*Data) {
	var dataMap = make(map[string]*Data)
	for _, d := range datum {
		dataMap[d.Data.GetJSONKey()] = d
	}

	text := http.StatusText(status)
	res := Response{
		Status:     &status,
		StatusDesc: &text,
		Data:       dataMap,
	}

	render.Status(r, status)
	render.SetContentType(render.ContentTypeJSON)
	render.JSON(w, r, res)
	return
}

// RespondOk ...
func RespondOk(w http.ResponseWriter, r *http.Request, datum ...*Data) {
	Respond(w, r, http.StatusOK, datum...)
	return
}

// RespondError ...
func RespondInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	RespondError(w, r, http.StatusInternalServerError, err)
	return
}

// RespondError ...
func RespondError(w http.ResponseWriter, r *http.Request, status int, err error) {
	errs := []string{err.Error()}
	text := http.StatusText(status)
	res := Response{
		Status:     &status,
		StatusDesc: &text,
		Errors:     &errs,
	}

	render.Status(r, status)
	render.SetContentType(render.ContentTypeJSON)
	render.JSON(w, r, res)
	return
}
