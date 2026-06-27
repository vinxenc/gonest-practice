package employeeModule

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// EmployeeReader is the data-access surface the service depends on. Defining it
// as an interface lets the service be tested with a fake, and lets the concrete
// GORM Repository be swapped without touching callers.
type EmployeeReader interface {
	List(ctx context.Context, limit, offset int) (employees []Employee, total int64, err error)
}

// Repository is the GORM-backed data-access layer for the employee module.
type Repository struct {
	db *gorm.DB
}

// EmployeeRepository constructs an employee Repository with its database handle
// (fx provider).
func EmployeeRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List returns a page of employees ordered by id, along with the total number of
// employees for pagination metadata.
func (r *Repository) List(ctx context.Context, limit, offset int) ([]Employee, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&Employee{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting employees: %w", err)
	}

	var employees []Employee
	if err := r.db.WithContext(ctx).
		Order("id").
		Limit(limit).
		Offset(offset).
		Find(&employees).Error; err != nil {
		return nil, 0, fmt.Errorf("listing employees: %w", err)
	}

	return employees, total, nil
}
