package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m1/go-service-healthcheck/model/mock"
)

func TestServiceMetric_UnmarshalJSON(t *testing.T) {
	metricData := map[string]interface{}{
		"id":             1,
		"service_id":     1,
		"tick_count":     5,
		"up_count":       5,
		"down_count":     0,
		"uptime_percent": 100,
		"created_at":     mock.Time.Unix(),
		"updated_at":     mock.Time.Unix(),
	}

	bytes, err := json.Marshal(&metricData)
	assert.NoError(t, err)

	var metric ServiceMetric
	err = json.Unmarshal(bytes, &metric)
	assert.NoError(t, err)

	assert.Equal(t, metric.CreatedAt.Unix(), metricData["created_at"])
	assert.Equal(t, metric.UpdatedAt.Unix(), metricData["updated_at"])
}
