CREATE TABLE payslip (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    period_id UUID NOT NULL REFERENCES attendance_period(id),
    base_salary NUMERIC NOT NULL,
    present_days INT NOT NULL,
    working_days INT NOT NULL,
    overtime_hours INT NOT NULL,
    reimbursement NUMERIC NOT NULL,
    take_home_pay NUMERIC NOT NULL,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    ip_address TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);