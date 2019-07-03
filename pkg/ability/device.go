package ability

// Device information
type Device struct {
	// Some dynamic information sent with each request
	State interface{} `json:"state,omitempty"`
}
