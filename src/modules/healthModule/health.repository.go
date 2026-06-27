package healthModule

// Repository is the data-access layer for the health module.
type Repository struct{}

// HealthRepository constructs a health Repository (fx provider).
func HealthRepository() *Repository {
	return &Repository{}
}

// Status reports the underlying health status.
func (r *Repository) Status() string {
	return "ok"
}
