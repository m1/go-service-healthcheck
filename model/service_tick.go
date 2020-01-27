package model

import (
	"encoding/json"
	"time"
)

type ServiceTick struct {
	Model
	ServiceID      int       `json:"service_id"`
	IsUp           bool      `json:"is_healthy"`
	ResponseTimeMS float64   `json:"response_time_ms"`
	CreatedAt      time.Time `json:"created_at"`
}

func (s *ServiceTick) UnmarshalJSON(data []byte) error {
	type Alias ServiceTick

	j := &struct {
		CreatedAt int64 `json:"created_at"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	s.CreatedAt = time.Unix(j.CreatedAt, 0).UTC()

	return nil
}
