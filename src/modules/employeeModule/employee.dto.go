package employeeModule

// dateLayout formats date-only fields (birth_date, hire_date) for JSON output.
const dateLayout = "2006-01-02"

// defaultLimit and maxLimit bound the page size for the list endpoint.
const (
	defaultLimit = 20
	maxLimit     = 100
)

// ListEmployeesInput is the query for GET /employees, with offset pagination.
type ListEmployeesInput struct {
	Limit  int `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Maximum number of employees to return"`
	Offset int `query:"offset" minimum:"0" default:"0" doc:"Number of employees to skip"`
}

// EmployeeDTO is the API representation of an employee.
type EmployeeDTO struct {
	ID        int64  `json:"id" doc:"Employee identifier" example:"10001"`
	BirthDate string `json:"birthDate" doc:"Date of birth (YYYY-MM-DD)" example:"1953-09-02"`
	FirstName string `json:"firstName" doc:"Given name" example:"Georgi"`
	LastName  string `json:"lastName" doc:"Family name" example:"Facello"`
	Gender    string `json:"gender" doc:"Gender (M or F)" example:"M"`
	HireDate  string `json:"hireDate" doc:"Date hired (YYYY-MM-DD)" example:"1986-06-26"`
}

// ListEmployeesOutput is the response body for GET /employees.
type ListEmployeesOutput struct {
	Body struct {
		Employees []EmployeeDTO `json:"employees" doc:"Page of employees"`
		Limit     int           `json:"limit" doc:"Page size that was applied"`
		Offset    int           `json:"offset" doc:"Offset that was applied"`
		Total     int64         `json:"total" doc:"Total number of employees"`
	}
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
