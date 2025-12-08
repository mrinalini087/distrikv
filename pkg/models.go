package pkg

// SetRequest is the data structure for saving a key-value pair.
type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}