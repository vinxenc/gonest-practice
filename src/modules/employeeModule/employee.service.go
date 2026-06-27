package employeeModule

import "context"

// Service holds the business logic for the employee module.
type Service struct {
	repo EmployeeReader
}

// EmployeeService constructs an employee Service with its repository (fx
// provider). It depends on the EmployeeReader interface so the data layer can be
// faked in tests.
func EmployeeService(repo EmployeeReader) *Service {
	return &Service{repo: repo}
}

// List returns a page of employees and the total count, normalizing the page
// size to a sane default and ceiling before delegating to the repository.
func (s *Service) List(ctx context.Context, limit, offset int) ([]Employee, int64, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}
