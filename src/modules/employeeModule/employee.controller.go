package employeeModule

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// Controller registers and handles the employee module's HTTP routes.
type Controller struct {
	service *Service
}

// EmployeeController constructs an employee Controller with its service (fx
// provider).
func EmployeeController(service *Service) *Controller {
	return &Controller{service: service}
}

// RegisterRoutes wires the employee endpoints onto the given Huma API.
func (c *Controller) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-employees",
		Method:      http.MethodGet,
		Path:        "/employees",
		Summary:     "List employees",
		Description: "Returns a paginated list of employees ordered by id.",
		Tags:        []string{"Employees"},
	}, c.list)
}

// list handles GET /employees.
func (c *Controller) list(ctx context.Context, input *ListEmployeesInput) (*ListEmployeesOutput, error) {
	employees, total, err := c.service.List(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to list employees", err)
	}

	resp := &ListEmployeesOutput{}
	resp.Body.Employees = toEmployeeDTOs(employees)
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset
	resp.Body.Total = total
	return resp, nil
}
