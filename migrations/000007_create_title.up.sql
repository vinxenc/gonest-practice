CREATE TABLE employees.title (
    employee_id bigint      NOT NULL,
    title       varchar(50) NOT NULL,
    from_date   date        NOT NULL,
    to_date     date,
    PRIMARY KEY (employee_id, title, from_date),
    CONSTRAINT fk_title_employee
        FOREIGN KEY (employee_id) REFERENCES employees.employee (id)
        ON UPDATE RESTRICT ON DELETE CASCADE
);
