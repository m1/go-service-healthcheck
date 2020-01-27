package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m1/go-service-healthcheck/model/mock"
)

func TestServiceTick_UnmarshalJSON(t *testing.T) {
	tickData := map[string]interface{}{
		"id":               1,
		"is_up":            1,
		"response_time_ms": 100,
		"created_at":       mock.Time.Unix(),
	}

	bytes, err := json.Marshal(&tickData)
	assert.NoError(t, err)

	var tick ServiceTick
	err = json.Unmarshal(bytes, &tick)
	assert.NoError(t, err)

	assert.Equal(t, tick.CreatedAt.Unix(), tickData["created_at"])
}
