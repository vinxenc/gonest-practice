package employeeModule

import "testing"

// TestEntities_TableName verifies each entity maps to its schema-qualified table
// in the employees schema.
func TestEntities_TableName(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Employee", Employee{}.TableName(), "employees.employee"},
		{"Department", Department{}.TableName(), "employees.department"},
		{"DepartmentEmployee", DepartmentEmployee{}.TableName(), "employees.department_employee"},
		{"DepartmentManager", DepartmentManager{}.TableName(), "employees.department_manager"},
		{"Salary", Salary{}.TableName(), "employees.salary"},
		{"Title", Title{}.TableName(), "employees.title"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("%s.TableName() = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}
