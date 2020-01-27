package runner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m1/go-service-healthcheck/log"
	"github.com/m1/go-service-healthcheck/model"
	"github.com/m1/go-service-healthcheck/repositories/mock"
)

func SetupTestServer() *httptest.Server {
	testSrv := http.NewServeMux()
	testSrv.HandleFunc("/Test_worker_queryService_valid_is_up", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `ok`)
	})
	testSrv.HandleFunc("/Test_worker_queryService_valid_is_down", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	srv := httptest.NewServer(testSrv)
	return srv
}

func NewMockRepositories() *Repositories {
	return &Repositories{Service: mock.NewServiceRepository(
		mock.Services,
		mock.Metrics,
		mock.Ticks,
		mock.Events,
	)}
}

func Test_worker_queryService(t *testing.T) {
	srv := SetupTestServer()
	type fields struct {
		Logger       *log.Logger
		Repositories *Repositories
		Client       http.Client
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		isUp   bool
	}{
		{
			name: "valid is up",
			fields: fields{
				Logger:       log.NewBasicLogger(),
				Repositories: NewMockRepositories(),
				Client:       http.Client{},
			},
			args: args{url: fmt.Sprintf("%v/Test_worker_queryService_valid_is_up", srv.URL)},
			isUp: true,
		},
		{
			name: "valid is down",
			fields: fields{
				Logger:       log.NewBasicLogger(),
				Repositories: NewMockRepositories(),
				Client:       http.Client{},
			},
			args: args{url: fmt.Sprintf("%v/Test_worker_queryService_valid_is_down", srv.URL)},
			isUp: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := worker{
				Logger:       tt.fields.Logger,
				Repositories: tt.fields.Repositories,
				Client:       tt.fields.Client,
			}
			got, _ := w.queryService(tt.args.url)
			if got != tt.isUp {
				t.Errorf("queryService() got = %v, want %v", got, tt.isUp)
			}
		})
	}
}

func Test_worker_calculateMovingAverage(t *testing.T) {
	type args struct {
		metric   model.ServiceMetric
		duration float64
	}
	tests := []struct {
		name   string
		args   args
		want   float64
	}{
		{
			name:   "valid",
			args:   args{
				metric:   model.ServiceMetric{
					TickCount:             5,
					AverageResponseTimeMS: 500,
				},
				duration: 50,
			},
			want:   425,
		},
		{
			name:   "valid",
			args:   args{
				metric:   model.ServiceMetric{
					TickCount:             5,
					AverageResponseTimeMS: 500,
				},
				duration: 80,
			},
			want:   430,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := worker{}
			if got := w.calculateMovingAverage(tt.args.metric, tt.args.duration); got != tt.want {
				t.Errorf("calculateMovingAverage() = %v, want %v", got, tt.want)
			}
		})
	}
}