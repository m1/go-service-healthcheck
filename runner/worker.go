package runner

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/m1/go-service-healthcheck/config"
	"github.com/m1/go-service-healthcheck/log"
	"github.com/m1/go-service-healthcheck/model"
)

type worker struct {
	Logger       *log.Logger
	Repositories *Repositories
	Client       http.Client
}

func newWorker(config *config.Config, logger *log.Logger, repositories *Repositories) *worker {
	worker := &worker{
		Logger:       logger,
		Repositories: repositories,
		Client: http.Client{
			Timeout: time.Duration(config.RunnerConfig.HTTPTimeoutSeconds) * time.Second,
		}}
	return worker
}

func (w worker) work(wg *sync.WaitGroup, service model.Service) {
	defer wg.Done()
	isUp, duration := w.queryService(service.URL)

	// get metric if has one
	w.updateMetric(service, duration, isUp)
	w.saveTick(service, duration, isUp)
	w.updateEvent(service, isUp)

	w.Logger.Info(fmt.Sprintf("fetched service: %v up=%v duration=%v", service.ID, isUp, duration))
}

func (w worker) queryService(url string) (bool, float64) {
	start := time.Now()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, 0
	}
	resp, err := w.Client.Do(req)
	durationMs := time.Since(start).Seconds() * 1000
	if err != nil {
		return false, durationMs
	}

	if resp.StatusCode >= 400 {
		return false, durationMs
	}

	return true, durationMs
}

func (w worker) calculateMovingAverage(metric model.ServiceMetric, duration float64) float64 {
	return ((metric.AverageResponseTimeMS * float64(metric.TickCount)) + duration) / float64(metric.TickCount+1)
}

func (w worker) updateMetric(service model.Service, duration float64, up bool) {
	metric, err := w.Repositories.Service.GetMetric(service.ID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			metric = model.ServiceMetric{
				ServiceID:             service.ID,
				TickCount:             1,
				AverageResponseTimeMS: duration,
			}
			if up {
				metric.UpCount = 1
			} else {
				metric.DownCount = 1
			}

			err := w.Repositories.Service.SaveMetric(&metric)
			if err != nil {
				w.Logger.With(zap.Error(err)).Error("err saving service metric")
			}
		} else {
			w.Logger.With(zap.Error(err)).Error("err getting service metric")
		}
		return
	}
	metric.TickCount++
	metric.AverageResponseTimeMS = w.calculateMovingAverage(metric, duration)
	if up {
		metric.UpCount++
	} else {
		metric.DownCount++
	}

	err = w.Repositories.Service.UpdateMetric(&metric)
	if err != nil {
		w.Logger.With(zap.Error(err)).Error("err updating service metric")
	}
}

func (w worker) saveTick(service model.Service, duration float64, up bool) {
	tick := model.ServiceTick{
		ServiceID:      service.ID,
		IsUp:           up,
		ResponseTimeMS: duration,
		CreatedAt:      time.Now(),
	}
	err := w.Repositories.Service.SaveTick(&tick)
	if err != nil {
		w.Logger.With(zap.Error(err)).Error("err saving service tick")
	}
}

func (w worker) updateEvent(service model.Service, up bool) {
	// get latest event
	event, err := w.Repositories.Service.GetLatestEvent(service.ID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			event := model.ServiceEvent{
				ServiceID:   service.ID,
				Event:       model.ServiceEventUptime,
				DateStarted: time.Now(),
			}
			if !up {
				event.Event = model.ServiceEventDowntime
			}
			err := w.Repositories.Service.SaveEvent(&event)
			if err != nil {
				w.Logger.With(zap.Error(err)).Error("err saving service event")
			}
		} else {
			w.Logger.With(zap.Error(err)).Error("err getting latest service event")
		}
		return
	}

	eventStr := model.ServiceEventUptime
	if !up {
		eventStr = model.ServiceEventDowntime
	}

	if eventStr != event.Event {
		// update old event - start new
		now := time.Now()
		event.DateEnded = &now
		err := w.Repositories.Service.UpdateEvent(&event)
		if err != nil {
			w.Logger.With(zap.Error(err)).Error("err updating old service event")
		}

		event := model.ServiceEvent{
			ServiceID:   service.ID,
			Event:       eventStr,
			DateStarted: time.Now(),
		}
		err = w.Repositories.Service.SaveEvent(&event)
		if err != nil {
			w.Logger.With(zap.Error(err)).Error("err saving service event")
		}
	}
}
