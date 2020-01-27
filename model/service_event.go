package model

import (
	"encoding/json"
	"time"
)

const (
	ServiceEventUptime   = "uptime"
	ServiceEventDowntime = "downtime"
)

type ServiceEvent struct {
	Model
	ID          int        `json:"id" gorm:"primary_key"`
	ServiceID   int        `json:"-"`
	Event       string     `json:"event"`
	DateStarted time.Time  `json:"date_started"`
	DateEnded   *time.Time `json:"date_ended"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (ServiceEvent) GetLatestEventJSONKey() string {
	return "latest_event"
}

func (s *ServiceEvent) MarshalJSON() ([]byte, error) {
	type Alias ServiceEvent
	j := &struct {
		DateStarted int64  `json:"date_started"`
		DateEnded   *int64 `json:"date_ended"`
		CreatedAt   int64  `json:"created_at"`
		UpdatedAt   int64  `json:"updated_at"`
		*Alias
	}{
		DateStarted: s.DateStarted.Unix(),
		CreatedAt:   s.CreatedAt.Unix(),
		UpdatedAt:   s.UpdatedAt.Unix(),
		Alias:       (*Alias)(s),
	}

	if s.DateEnded != nil {
		unix := s.DateEnded.Unix()
		j.DateEnded = &unix
	}

	return json.Marshal(s.GetJSONMap(j))
}

func (s *ServiceEvent) UnmarshalJSON(data []byte) error {
	type Alias ServiceEvent

	j := &struct {
		DateStarted int64  `json:"date_started"`
		DateEnded   *int64 `json:"date_ended"`
		CreatedAt   int64  `json:"created_at"`
		UpdatedAt   int64  `json:"updated_at"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	s.DateStarted = time.Unix(j.DateStarted, 0).UTC()
	s.DateEnded = nil
	if j.DateEnded != nil {
		u := time.Unix(*j.DateEnded, 0).UTC()
		s.DateEnded = &u
	}
	s.CreatedAt = time.Unix(j.CreatedAt, 0).UTC()
	s.UpdatedAt = time.Unix(j.UpdatedAt, 0).UTC()

	return nil
}

type ServiceEvents []ServiceEvent

func (ServiceEvents) GetJSONKey() string {
	return "events"
}
