package healthModule

// HealthOutput is the response body for the health check endpoint.
type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"Service health status"`
	}
}
