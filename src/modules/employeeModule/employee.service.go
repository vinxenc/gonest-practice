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

// ListResult is a page of employees together with the total row count and the
// pagination that was actually applied (after clamping), so callers can report
// the effective limit/offset rather than the raw request values.
type ListResult struct {
	Employees []Employee
	Total     int64
	Limit     int
	Offset    int
}

// List returns a page of employees and the total count, normalizing the page
// size to a sane default and ceiling before delegating to the repository. The
// returned ListResult carries the normalized limit/offset that were applied.
func (s *Service) List(ctx context.Context, limit, offset int) (ListResult, error) {
	limit, offset = normalize(limit, offset)

	employees, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return ListResult{}, err
	}

	return ListResult{Employees: employees, Total: total, Limit: limit, Offset: offset}, nil
}

// normalize clamps the page size to [1, maxLimit] (defaulting an unset size) and
// floors the offset at zero.
func normalize(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
