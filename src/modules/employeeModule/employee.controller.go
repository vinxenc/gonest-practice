package employeeModule

import (
	"net/http"
	"strconv"

	"github.com/0xfurai/gonest"
)

// Controller registers and handles the employee module's HTTP routes.
type Controller struct {
	service *Service
	logger  gonest.Logger
}

// EmployeeController constructs an employee Controller with its service and the
// framework logger (both gonest providers). The logger lets the handler record
// the underlying cause of a failure server-side while still returning a generic
// error to the client.
func EmployeeController(service *Service, logger gonest.Logger) *Controller {
	return &Controller{service: service, logger: logger}
}

// Register wires the employee endpoints onto the given router and declares their
// OpenAPI metadata (summary, tag, response schema) for the Swagger document.
func (c *Controller) Register(r gonest.Router) {
	r.Get("/employees", c.list).
		Summary("List employees").
		Tags("Employees").
		Response(http.StatusOK, ListEmployeesResponse{})
}

// list handles GET /employees. The limit/offset query parameters are advisory:
// the service is the single source of truth for pagination bounds and clamps
// out-of-range values rather than rejecting them, so the response echoes the
// effective limit/offset that were actually applied.
func (c *Controller) list(ctx gonest.Context) error {
	limit := queryInt(ctx, "limit")
	offset := queryInt(ctx, "offset")

	result, err := c.service.List(ctx.Ctx(), limit, offset)
	if err != nil {
		// Log the root cause for observability, but return a generic message so
		// internal error details are never exposed to the client.
		c.logger.Error("failed to list employees: %v", err)
		return gonest.NewInternalServerError("failed to list employees")
	}

	return ctx.JSON(http.StatusOK, ListEmployeesResponse{
		Employees: toEmployeeDTOs(result.Employees),
		Limit:     result.Limit,
		Offset:    result.Offset,
		Total:     result.Total,
	})
}

// queryInt reads a query parameter as an int, returning 0 when it is absent or
// not a valid integer. The service clamps 0 to its default, so a missing or
// malformed value is treated as "unspecified" rather than an error — preserving
// the endpoint's clamp-don't-reject contract.
func queryInt(ctx gonest.Context, name string) int {
	v, err := strconv.Atoi(ctx.Query(name))
	if err != nil {
		return 0
	}
	return v
}
