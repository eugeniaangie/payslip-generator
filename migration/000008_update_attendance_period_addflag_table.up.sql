ALTER TABLE attendance_period
ADD COLUMN is_payroll_processed BOOLEAN NOT NULL DEFAULT FALSE;
