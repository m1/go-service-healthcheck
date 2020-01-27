package model

import (
	"encoding/json"
	"time"
)

type Service struct {
	Model
	ID        int        `json:"id" gorm:"primary_key"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	Metric *ServiceMetric `json:"metric" gorm:"foreignkey:ID;association_foreignkey:ServiceID"`

	CurrentStatusID *int          `json:"-" gorm:"-"`
	CurrentStatus   *ServiceEvent `json:"current_status" gorm:"foreignkey:CurrentStatusID;association_foreignkey:ID"`
}

func (s Service) GetJSONKey() string {
	return "service"
}

func (s *Service) MarshalJSON() ([]byte, error) {
	type Alias Service
	j := s.GetJSONMap(&struct {
		CreatedAt int64 `json:"created_at"`
		UpdatedAt int64 `json:"updated_at"`
		*Alias
	}{
		CreatedAt: s.CreatedAt.Unix(),
		UpdatedAt: s.UpdatedAt.Unix(),
		Alias:     (*Alias)(s),
	})

	if s.CurrentStatus != nil {
		j[s.CurrentStatus.GetLatestEventJSONKey()] = NestedModel{Data: &s.CurrentStatus}
	}

	if s.Metric != nil {
		j[s.Metric.GetJSONKey()] = NestedModel{Data: &s.Metric}
	}

	return json.Marshal(j)
}

func (s *Service) UnmarshalJSON(data []byte) error {
	type Alias Service

	j := &struct {
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64 `json:"updated_at"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	s.CreatedAt = time.Unix(j.CreatedAt, 0).UTC()
	s.UpdatedAt = time.Unix(j.UpdatedAt, 0).UTC()

	return nil
}

type Services []Service

func (Services) GetJSONKey() string {
	return "services"
}
