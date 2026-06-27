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
