package employeeModule

// dateLayout formats date-only fields (birth_date, hire_date) for JSON output.
const dateLayout = "2006-01-02"

// defaultLimit and maxLimit bound the page size for the list endpoint.
const (
	defaultLimit = 20
	maxLimit     = 100
)

// EmployeeDTO is the API representation of an employee.
type EmployeeDTO struct {
	ID        int64  `json:"id"`
	BirthDate string `json:"birthDate"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	HireDate  string `json:"hireDate"`
}

// ListEmployeesResponse is the response body for GET /employees: a page of
// employees together with the pagination that was actually applied (after the
// service clamps the requested values) and the total row count.
type ListEmployeesResponse struct {
	Employees []EmployeeDTO `json:"employees"`
	Limit     int           `json:"limit"`
	Offset    int           `json:"offset"`
	Total     int64         `json:"total"`
}

// toEmployeeDTO converts a persisted Employee into its API representation.
func toEmployeeDTO(e Employee) EmployeeDTO {
	return EmployeeDTO{
		ID:        e.ID,
		BirthDate: e.BirthDate.Format(dateLayout),
		FirstName: e.FirstName,
		LastName:  e.LastName,
		Gender:    e.Gender,
		HireDate:  e.HireDate.Format(dateLayout),
	}
}

// toEmployeeDTOs converts a slice of employees, always returning a non-nil slice
// so the JSON response is `[]` rather than `null` when empty.
func toEmployeeDTOs(employees []Employee) []EmployeeDTO {
	dtos := make([]EmployeeDTO, 0, len(employees))
	for _, e := range employees {
		dtos = append(dtos, toEmployeeDTO(e))
	}
	return dtos
}
