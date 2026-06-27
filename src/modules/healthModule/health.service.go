package healthModule

// Service holds the business logic for the health module.
type Service struct {
	repo *Repository
}

// HealthService constructs a health Service with its repository (fx provider).
func HealthService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Check returns the current health status of the service.
func (s *Service) Check() string {
	return s.repo.Status()
}
