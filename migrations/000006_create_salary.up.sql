CREATE TABLE employees.salary (
    employee_id bigint NOT NULL,
    amount      bigint NOT NULL,
    from_date   date   NOT NULL,
    to_date     date   NOT NULL,
    PRIMARY KEY (employee_id, from_date),
    CONSTRAINT fk_salary_employee
        FOREIGN KEY (employee_id) REFERENCES employees.employee (id)
        ON UPDATE RESTRICT ON DELETE CASCADE
);
