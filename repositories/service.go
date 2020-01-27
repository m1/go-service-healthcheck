package repositories

import (
	"github.com/jinzhu/gorm"

	"github.com/m1/go-service-healthcheck/model"
	"github.com/m1/go-service-healthcheck/storage"
)

type ServiceRepository interface {
	Get(id int) (model.Service, error)
	GetAll() (model.Services, error)
	GetMetric(id int) (model.ServiceMetric, error)
	SaveMetric(metric *model.ServiceMetric) error
	UpdateMetric(metric *model.ServiceMetric) error
	SaveTick(tick *model.ServiceTick) error
	GetLatestEvent(id int) (model.ServiceEvent, error)
	SaveEvent(event *model.ServiceEvent) error
	UpdateEvent(event *model.ServiceEvent) error
	GetEvents(id int) (model.ServiceEvents, error)
}

type serviceRepository struct {
	*storage.SQLiteDB
}

func NewServiceRepository(db *storage.SQLiteDB) ServiceRepository {
	return &serviceRepository{SQLiteDB: db}
}

func (r serviceRepository) Get(id int) (model.Service, error) {
	service := model.Service{ID: id}
	err := r.DB.Find(&service).Error
	return service, err
}

func (r serviceRepository) GetAll() (model.Services, error) {
	var services model.Services
	err := r.DB.Select([]string{
		"services.*",
		"service_events.id AS current_status_id",
	}).Scopes(r.scopeGetCurrentStatusID()).
		Preload("Metric").
		Find(&services).
		Error
	return services, err
}

func (r serviceRepository) scopeGetCurrentStatusID() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(`
			LEFT JOIN service_events ON 
				service_events.id = ( 
					SELECT latest.id
					FROM service_events AS latest
					WHERE latest.service_id = services.id
					ORDER BY latest.id DESC
					LIMIT 1
				)
		`).Preload("CurrentStatus")
	}
}

func (r serviceRepository) GetMetric(id int) (model.ServiceMetric, error) {
	metric := model.ServiceMetric{ServiceID: id}
	err := r.DB.Find(&metric).Error
	return metric, err
}

func (r serviceRepository) SaveMetric(metric *model.ServiceMetric) error {
	r.DB.NewRecord(&metric)
	return r.Create(&metric).Error
}

func (r serviceRepository) UpdateMetric(metric *model.ServiceMetric) error {
	return r.DB.Save(&metric).Error
}

func (r serviceRepository) SaveTick(tick *model.ServiceTick) error {
	r.DB.NewRecord(&tick)
	return r.Create(&tick).Error
}

func (r serviceRepository) GetLatestEvent(id int) (model.ServiceEvent, error) {
	event := model.ServiceEvent{}
	err := r.DB.Where("service_id = ?", id).
		Order("created_at desc").
		Limit(1).
		Find(&event).
		Error
	return event, err
}

func (r serviceRepository) GetEvents(id int) (model.ServiceEvents, error) {
	var events model.ServiceEvents
	err := r.DB.Where("service_id = ?", id).
		Order("created_at desc").
		Find(&events).Error
	return events, err
}

func (r serviceRepository) SaveEvent(event *model.ServiceEvent) error {
	r.DB.NewRecord(&event)
	return r.Create(&event).Error
}

func (r serviceRepository) UpdateEvent(event *model.ServiceEvent) error {
	return r.DB.Save(&event).Error
}
