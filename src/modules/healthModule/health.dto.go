package healthModule

// HealthResponse is the response body for the health check endpoint. The swagger
// tag supplies the example shown in the generated OpenAPI document.
type HealthResponse struct {
	Status string `json:"status" swagger:"example=ok"`
}
