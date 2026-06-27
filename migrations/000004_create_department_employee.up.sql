CREATE TABLE employees.department_employee (
    employee_id   bigint  NOT NULL,
    department_id char(4) NOT NULL,
    from_date     date    NOT NULL,
    to_date       date    NOT NULL,
    PRIMARY KEY (employee_id, department_id),
    CONSTRAINT fk_department_employee_employee
        FOREIGN KEY (employee_id) REFERENCES employees.employee (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    CONSTRAINT fk_department_employee_department
        FOREIGN KEY (department_id) REFERENCES employees.department (id)
        ON UPDATE RESTRICT ON DELETE CASCADE
);

-- The composite PK is keyed by employee_id first, so it can't serve
-- department-driven lookups or the department_id FK check; index it separately.
CREATE INDEX idx_department_employee_department_id
    ON employees.department_employee (department_id);
