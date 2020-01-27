package response

// Model is the model interface that all models must
// implement to return json data
type Model interface {
	// GetJSONKey returns the json key that the model will
	// use
	GetJSONKey() string
}
