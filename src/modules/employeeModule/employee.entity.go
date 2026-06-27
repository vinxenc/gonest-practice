package employeeModule

import "time"

// The neon "employees" sample database keeps every object in a dedicated
// `employees` schema (not `public`), so each entity's TableName is
// schema-qualified. See:
// https://github.com/neondatabase/postgres-sample-dbs#employees-database

// Employee maps the employees.employee table — one row per employee.
type Employee struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	BirthDate time.Time `gorm:"column:birth_date;type:date;not null"`
	FirstName string    `gorm:"column:first_name;not null"`
	LastName  string    `gorm:"column:last_name;not null"`
	Gender    string    `gorm:"column:gender;type:employees.employee_gender;not null"`
	HireDate  time.Time `gorm:"column:hire_date;type:date;not null"`

	// Associations (loaded only when explicitly Preloaded).
	DepartmentAssignments []DepartmentEmployee `gorm:"foreignKey:EmployeeID"`
	Salaries              []Salary             `gorm:"foreignKey:EmployeeID"`
	Titles                []Title              `gorm:"foreignKey:EmployeeID"`
}

// TableName returns the schema-qualified table for Employee.
func (Employee) TableName() string { return "employees.employee" }

// Department maps the employees.department table.
type Department struct {
	ID       string `gorm:"column:id;type:char(4);primaryKey"`
	DeptName string `gorm:"column:dept_name;not null"`
}

// TableName returns the schema-qualified table for Department.
func (Department) TableName() string { return "employees.department" }

// DepartmentEmployee maps employees.department_employee, the membership of an
// employee in a department over a time range (composite primary key).
type DepartmentEmployee struct {
	EmployeeID   int64     `gorm:"column:employee_id;primaryKey"`
	DepartmentID string    `gorm:"column:department_id;type:char(4);primaryKey"`
	FromDate     time.Time `gorm:"column:from_date;type:date;not null"`
	ToDate       time.Time `gorm:"column:to_date;type:date;not null"`
}

// TableName returns the schema-qualified table for DepartmentEmployee.
func (DepartmentEmployee) TableName() string { return "employees.department_employee" }

// DepartmentManager maps employees.department_manager, the management of a
// department by an employee over a time range (composite primary key).
type DepartmentManager struct {
	EmployeeID   int64     `gorm:"column:employee_id;primaryKey"`
	DepartmentID string    `gorm:"column:department_id;type:char(4);primaryKey"`
	FromDate     time.Time `gorm:"column:from_date;type:date;not null"`
	ToDate       time.Time `gorm:"column:to_date;type:date;not null"`
}

// TableName returns the schema-qualified table for DepartmentManager.
func (DepartmentManager) TableName() string { return "employees.department_manager" }

// Salary maps employees.salary, an employee's salary over a time range
// (composite primary key of employee_id + from_date).
type Salary struct {
	EmployeeID int64     `gorm:"column:employee_id;primaryKey"`
	Amount     int64     `gorm:"column:amount;not null"`
	FromDate   time.Time `gorm:"column:from_date;type:date;primaryKey"`
	ToDate     time.Time `gorm:"column:to_date;type:date;not null"`
}

// TableName returns the schema-qualified table for Salary.
func (Salary) TableName() string { return "employees.salary" }

// Title maps employees.title, a job title held by an employee over a time range
// (composite primary key of employee_id + title + from_date). ToDate is nullable
// for a currently-held title.
type Title struct {
	EmployeeID int64      `gorm:"column:employee_id;primaryKey"`
	Title      string     `gorm:"column:title;primaryKey"`
	FromDate   time.Time  `gorm:"column:from_date;type:date;primaryKey"`
	ToDate     *time.Time `gorm:"column:to_date;type:date"`
}

// TableName returns the schema-qualified table for Title.
func (Title) TableName() string { return "employees.title" }
