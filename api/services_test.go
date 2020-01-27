package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/m1/go-service-healthcheck/model"
	"github.com/m1/go-service-healthcheck/repositories/mock"
	"github.com/m1/go-service-healthcheck/response"
)

type ServicesTestSuite struct {
	T *testing.T

	API          *API
	Repositories *Repositories
	Handler      *ServicesHandler
}

func NewServicesTestSuite(t *testing.T) ServicesTestSuite {
	repos := &Repositories{Service: mock.NewServiceRepository(
		mock.Services,
		mock.Metrics,
		mock.Ticks,
		mock.Events,
	)}

	api := &API{router: chi.NewRouter(), Repositories: repos}

	handler := NewServicesHandler(api)
	suite := ServicesTestSuite{
		T:            t,
		API:          api,
		Repositories: repos,
		Handler:      handler,
	}
	suite.API.router.Mount("/", suite.Handler.GetRoutes())
	return suite
}

func (s *ServicesTestSuite) MakeRequest(r ServiceTestSuiteRequest, data interface{}) *httptest.ResponseRecorder {
	var body io.Reader
	if r.body != nil {
		body = bytes.NewReader(*r.body)
	}
	request, _ := http.NewRequest(r.method, r.url, body)
	request.Header.Set("Content-Type", "application/json")

	if r.params != nil {
		q := request.URL.Query()
		for k, v := range r.params {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}

	resp := httptest.NewRecorder()
	s.API.router.ServeHTTP(resp, request)
	err := json.Unmarshal(resp.Body.Bytes(), &data)
	if err != nil {
		s.T.Error(s.T, err)
	}
	return resp
}

type ServiceTestSuiteRequest struct {
	method string
	url    string
	body   *[]byte
	params map[string]string
}

func NewServiceTestSuiteRequest(method string, url string, body *[]byte, params map[string]string) *ServiceTestSuiteRequest {
	return &ServiceTestSuiteRequest{method: method, url: url, body: body, params: params}
}

func TestServicesHandler_GetServices(t *testing.T) {
	suite := NewServicesTestSuite(t)

	var res struct {
		response.Response
		Data struct {
			Services struct {
				Data model.Services `json:"data"`
			} `json:"services"`
		} `json:"data"`
	}
	recorder := suite.MakeRequest(
		ServiceTestSuiteRequest{
			method: http.MethodGet,
			url:    "/"},
		&res)

	assert.Equal(t, http.StatusOK, recorder.Code, "should be equal")
	assert.Equal(t, len(res.Data.Services.Data), len(mock.Services), "should be equal")
	assert.Equal(t, res.Data.Services.Data[0].CreatedAt, mock.Time, "should be equal")

	// should have service id 1, 2 and 3
	for _, service := range res.Data.Services.Data {
		if service.ID != 1 && service.ID != 2 && service.ID != 3 {
			t.Errorf("%v should not exist", service.ID)
		}
	}
}

func TestServicesHandler_GetEvents(t *testing.T) {
	suite := NewServicesTestSuite(t)

	type res struct {
		response.Response
		Data struct {
			Events struct {
				Data model.ServiceEvents `json:"data"`
			} `json:"events"`
		} `json:"data"`
	}

	tests := []struct {
		name string
		url  string
		test func(t *testing.T, res *res)
	}{
		{
			name: "valid service 1",
			url:  "/1/events",
			test: func(t *testing.T, res *res) {
				assert.Equal(t, http.StatusOK, *res.Status, "should be equal")
				assert.Equal(t, len(res.Data.Events.Data), 2, "should be equal")
			},
		},
		{
			name: "valid service 2",
			url:  "/2/events",
			test: func(t *testing.T, res *res) {
				assert.Equal(t, http.StatusOK, *res.Status, "should be equal")
				assert.Equal(t, len(res.Data.Events.Data), 1, "should be equal")
			},
		},
		{
			name: "invalid 404",
			url:  "/10/events",
			test: func(t *testing.T, res *res) {
				errs := *res.Errors
				assert.Equal(t, http.StatusNotFound, *res.Status, "should be equal")
				assert.Equal(t, errs[0], gorm.ErrRecordNotFound.Error(), "should be equal")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &res{}
			_ = suite.MakeRequest(
				ServiceTestSuiteRequest{
					method: http.MethodGet,
					url:    tt.url,
				},
				&res)
			tt.test(t, res)
		})
	}
}
