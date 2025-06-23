-- Admin Insert (password: admin123)
INSERT INTO users (id, username, password_hash, full_name, user_role, created_at, updated_at)
VALUES (
    'a1e02347-b56d-4076-b8f3-339c64f87fe4',
    'admin',
    '$2a$10$a8bdO0cpanNmf2DhYqTv8OwIzAInLj41JYq0RezR9kVYjAcqVvf36',
    'Admin',
    'admin',
    now(),
    now()
);

-- Step 4 & 5: Fake Employees + Salary
DO $$
DECLARE
    i INTEGER;
    user_id UUID;
    gaji INTEGER;
BEGIN
    FOR i IN 1..100 LOOP
        -- Insert employee
        INSERT INTO users (id, username, password_hash, full_name, user_role, created_at, updated_at)
        VALUES (
            uuid_generate_v4(),
            'employee_' || i,
            '$2a$10$7dZDKKT6smIu.5aZ8qv6ze0LS36.NDZ8bYUb1KwB5pjX77eYF3GJS', -- password: employee123
            'Employee ' || i,
            'employee',
            now(),
            now()
        )
        RETURNING id INTO user_id;

        -- Random salary between 3jt - 7jt
        gaji := floor(random() * (7000000 - 3000000 + 1)) + 3000000;

        -- Insert salary
        INSERT INTO salary (id, user_id, amount, is_active, created_by, updated_by, created_at, updated_at)
        VALUES (
            uuid_generate_v4(),
            user_id,
            gaji,
            true,
            'a1e02347-b56d-4076-b8f3-339c64f87fe4', -- admin
            'a1e02347-b56d-4076-b8f3-339c64f87fe4',
            now(),
            now()
        );
    END LOOP;
END $$;