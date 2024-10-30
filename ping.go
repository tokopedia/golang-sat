package sat

// PingResponse contains health check payload
type PingResponse struct {
	Buildhash string `json:"buildhash"`
	Sandbox   bool   `json:"sandbox"`
	Status    string `json:"status"`
}
