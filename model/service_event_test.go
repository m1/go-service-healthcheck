package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m1/go-service-healthcheck/model/mock"
)

func TestServiceEvent_UnmarshalJSON(t *testing.T) {
	eventData := map[string]interface{}{
		"id":           1,
		"service_id":   1,
		"event":        ServiceEventUptime,
		"date_started": mock.Time.Unix(),
		"date_ended":   mock.Time.Unix(),
		"created_at":   mock.Time.Unix(),
		"updated_at":   mock.Time.Unix(),
		"deleted_at":   nil,
	}

	bytes, err := json.Marshal(&eventData)
	assert.NoError(t, err)

	var event ServiceEvent
	err = json.Unmarshal(bytes, &event)
	assert.NoError(t, err)

	assert.Equal(t, event.CreatedAt.Unix(), eventData["created_at"])
	assert.Equal(t, event.UpdatedAt.Unix(), eventData["updated_at"])
}
