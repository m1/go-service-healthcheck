package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m1/go-service-healthcheck/model/mock"
)

func TestService_UnmarshalJSON(t *testing.T) {
	serviceData := map[string]interface{}{
		"id":         1,
		"name":       "mock-1",
		"url":        "localhost/1",
		"created_at": mock.Time.Unix(),
		"updated_at": mock.Time.Unix(),
		"deleted_at": nil,
	}

	bytes, err := json.Marshal(&serviceData)
	assert.NoError(t, err)

	var service Service
	err = json.Unmarshal(bytes, &service)
	assert.NoError(t, err)
	assert.Equal(t, service.CreatedAt.Unix(), serviceData["created_at"])
	assert.Equal(t, service.UpdatedAt.Unix(), serviceData["updated_at"])
}
