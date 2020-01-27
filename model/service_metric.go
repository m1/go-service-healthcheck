package model

import (
	"encoding/json"
	"time"
)

type ServiceMetric struct {
	Model
	ServiceID             int       `json:"-" gorm:"primary_key"`
	TickCount             int       `json:"tick_count"`
	UpCount               int       `json:"up_count"`
	DownCount             int       `json:"down_count"`
	AverageResponseTimeMS float64   `json:"average_response_time_ms"`
	UptimePercent         float64   `json:"uptime_percent" gorm:"-"`
	CreatedAt             time.Time `json:"-"`
	UpdatedAt             time.Time `json:"-"`
}

func (s *ServiceMetric) MarshalJSON() ([]byte, error) {
	type Alias ServiceMetric

	j := &struct {
		CreatedAt int64 `json:"created_at"`
		UpdatedAt int64 `json:"updated_at"`
		*Alias
	}{
		CreatedAt: s.CreatedAt.Unix(),
		UpdatedAt: s.UpdatedAt.Unix(),
		Alias:     (*Alias)(s),
	}

	j.UptimePercent = (float64(s.UpCount) / float64(s.TickCount)) * 100

	return json.Marshal(s.GetJSONMap(j))
}

func (s *ServiceMetric) GetJSONKey() string {
	return "metric"
}
