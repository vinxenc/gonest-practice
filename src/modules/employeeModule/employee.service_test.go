package employeeModule

import (
	"context"
	"errors"
	"testing"
)

// fakeReader is a test double for EmployeeReader that records the arguments it
// was called with and returns canned results.
type fakeReader struct {
	employees []Employee
	total     int64
	err       error

	gotLimit  int
	gotOffset int
}

func (f *fakeReader) List(_ context.Context, limit, offset int) ([]Employee, int64, error) {
	f.gotLimit = limit
	f.gotOffset = offset
	return f.employees, f.total, f.err
}

// TestService_List_Delegates verifies the service passes through to the
// repository and returns its employees and total.
func TestService_List_Delegates(t *testing.T) {
	repo := &fakeReader{
		employees: []Employee{{ID: 1}, {ID: 2}},
		total:     2,
	}
	svc := EmployeeService(repo)

	result, err := svc.List(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if result.Total != 2 || len(result.Employees) != 2 {
		t.Fatalf("List = (%d employees, total %d), want (2, 2)", len(result.Employees), result.Total)
	}
	if result.Limit != 10 || result.Offset != 5 {
		t.Fatalf("result pagination = limit %d / offset %d, want 10 / 5", result.Limit, result.Offset)
	}
	if repo.gotLimit != 10 || repo.gotOffset != 5 {
		t.Fatalf("repo called with limit=%d offset=%d, want 10/5", repo.gotLimit, repo.gotOffset)
	}
}

// TestService_List_Normalizes verifies the service clamps the page size and
// offset to sane values before delegating.
func TestService_List_Normalizes(t *testing.T) {
	tests := []struct {
		name                  string
		limit, offset         int
		wantLimit, wantOffset int
	}{
		{"zero limit -> default", 0, 0, defaultLimit, 0},
		{"negative limit -> default", -5, 0, defaultLimit, 0},
		{"over max limit -> max", 1000, 0, maxLimit, 0},
		{"negative offset -> zero", 20, -10, 20, 0},
		{"in range passthrough", 50, 100, 50, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeReader{}
			svc := EmployeeService(repo)

			result, err := svc.List(context.Background(), tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("List returned error: %v", err)
			}
			if repo.gotLimit != tt.wantLimit || repo.gotOffset != tt.wantOffset {
				t.Fatalf("repo called with limit=%d offset=%d, want %d/%d",
					repo.gotLimit, repo.gotOffset, tt.wantLimit, tt.wantOffset)
			}
			if result.Limit != tt.wantLimit || result.Offset != tt.wantOffset {
				t.Fatalf("result pagination = limit %d / offset %d, want %d/%d",
					result.Limit, result.Offset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}

// TestService_List_PropagatesError verifies a repository error reaches the
// caller unchanged.
func TestService_List_PropagatesError(t *testing.T) {
	wantErr := errors.New("db down")
	svc := EmployeeService(&fakeReader{err: wantErr})

	if _, err := svc.List(context.Background(), 20, 0); !errors.Is(err, wantErr) {
		t.Fatalf("List error = %v, want %v", err, wantErr)
	}
}
