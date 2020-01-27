package mock

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/m1/go-service-healthcheck/model"
)

var Time = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
var Time2 = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

var Services = map[int]model.Service{
	1: {
		ID:        1,
		Name:      "mock-1",
		URL:       "localhost/mock-1",
		CreatedAt: Time,
		UpdatedAt: Time,
		DeletedAt: nil,
	},
	2: {
		ID:        2,
		Name:      "mock-2",
		URL:       "localhost/mock-2",
		CreatedAt: time.Unix(1577836800, 0),
		UpdatedAt: time.Unix(1577836800, 0),
		DeletedAt: nil,
	},
	3: {
		ID:        3,
		Name:      "mock-3",
		URL:       "localhost/mock-3",
		CreatedAt: time.Unix(1577836800, 0),
		UpdatedAt: time.Unix(1577836800, 0),
		DeletedAt: nil,
	},
}

var Metrics = map[int]model.ServiceMetric{
	1: {
		ServiceID:             1,
		TickCount:             2,
		UpCount:               1,
		DownCount:             1,
		AverageResponseTimeMS: 100,
		CreatedAt:             Time,
		UpdatedAt:             Time,
	},
	2: {
		ServiceID:             2,
		TickCount:             2,
		UpCount:               2,
		DownCount:             0,
		AverageResponseTimeMS: 200,
		CreatedAt:             Time,
		UpdatedAt:             Time,
	},
	3: {
		ServiceID:             3,
		TickCount:             4,
		UpCount:               0,
		DownCount:             4,
		AverageResponseTimeMS: 400,
		CreatedAt:             Time,
		UpdatedAt:             Time,
	},
}

var Ticks []model.ServiceTick

var Events = model.ServiceEvents{
	{
		ID:          1,
		ServiceID:   1,
		Event:       model.ServiceEventUptime,
		DateStarted: Time,
		DateEnded:   &Time2,
		CreatedAt:   Time,
		UpdatedAt:   Time2,
	},
	{
		ID:          2,
		ServiceID:   1,
		Event:       model.ServiceEventUptime,
		DateStarted: Time2,
		CreatedAt:   Time,
		UpdatedAt:   Time,
	},
	{
		ID:          3,
		ServiceID:   2,
		Event:       model.ServiceEventUptime,
		DateStarted: Time2,
		CreatedAt:   Time,
		UpdatedAt:   Time,
	},
}

type serviceRepository struct {
	MockServices map[int]model.Service
	MockMetrics  map[int]model.ServiceMetric
	MockTicks    []model.ServiceTick
	MockEvents   model.ServiceEvents
}

func NewServiceRepository(services map[int]model.Service, metrics map[int]model.ServiceMetric, ticks []model.ServiceTick, events model.ServiceEvents) *serviceRepository {
	return &serviceRepository{MockServices: services, MockMetrics: metrics, MockTicks: ticks, MockEvents: events}
}

func (s *serviceRepository) Get(id int) (model.Service, error) {
	service, ok := s.MockServices[id]
	if !ok {
		return model.Service{}, gorm.ErrRecordNotFound
	}
	return service, nil
}

func (s *serviceRepository) GetAll() (model.Services, error) {
	var services model.Services
	for _, s := range s.MockServices {
		services = append(services, s)
	}
	return services, nil
}

func (s *serviceRepository) GetMetric(id int) (model.ServiceMetric, error) {
	metric, ok := s.MockMetrics[id]
	if !ok {
		return model.ServiceMetric{}, gorm.ErrRecordNotFound
	}
	return metric, nil
}

func (s *serviceRepository) SaveMetric(metric *model.ServiceMetric) error {
	s.MockMetrics[metric.ServiceID] = *metric
	return nil
}

func (s *serviceRepository) UpdateMetric(metric *model.ServiceMetric) error {
	s.MockMetrics[metric.ServiceID] = *metric
	return nil
}

func (s *serviceRepository) SaveTick(tick *model.ServiceTick) error {
	s.MockTicks = append(s.MockTicks, *tick)
	return nil
}

func (s *serviceRepository) GetLatestEvent(id int) (model.ServiceEvent, error) {
	var latestEvent model.ServiceEvent
	for _, event := range s.MockEvents {
		if event.ServiceID == id && event.ID > latestEvent.ID {
			latestEvent = event
		}
	}
	return latestEvent, nil
}

func (s *serviceRepository) SaveEvent(event *model.ServiceEvent) error {
	s.MockEvents = append(s.MockEvents, *event)
	return nil
}

func (s serviceRepository) UpdateEvent(event *model.ServiceEvent) error {
	for i, mockEvent := range s.MockEvents {
		if event.ID == mockEvent.ID {
			s.MockEvents[i] = *event
		}
	}
	return nil
}

func (s serviceRepository) GetEvents(id int) (model.ServiceEvents, error) {
	_, ok := s.MockServices[id]
	if !ok {
		return model.ServiceEvents{}, gorm.ErrRecordNotFound
	}

	var events model.ServiceEvents
	for _, event := range s.MockEvents {
		if event.ServiceID == id {
			events = append(events, event)
		}
	}
	return events, nil
}
