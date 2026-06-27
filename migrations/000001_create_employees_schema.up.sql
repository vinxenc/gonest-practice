-- Prerequisite for every employees.* table: the dedicated schema and the
-- gender enum type used by employees.employee. All other migrations create
-- objects inside this schema.
CREATE SCHEMA IF NOT EXISTS employees;

CREATE TYPE employees.employee_gender AS ENUM ('M', 'F');
