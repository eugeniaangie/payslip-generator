CREATE TABLE reimbursement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    date DATE NOT NULL,
    amount NUMERIC NOT NULL,
    description TEXT,
    created_by UUID,
    updated_by UUID,
    ip_address TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
