package model

import (
	"encoding/json"
)

type Model struct {}

// GetJSONMap creates a map[string]interface from the
// passed interface, enables you to add mapped keys to the
// interface
func (Model) GetJSONMap(i interface{}) map[string]interface{} {
	m, _ := json.Marshal(i)
	var a interface{}
	_ = json.Unmarshal(m, &a)
	return a.(map[string]interface{})
}

type NestedModel struct {
	Data interface{} `json:"data"`
}
