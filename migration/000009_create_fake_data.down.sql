-- Down: Delete all employees and their salaries
DO $$
BEGIN
    DELETE FROM salary
    WHERE user_id IN (SELECT id FROM users WHERE user_role = 'employee');

    DELETE FROM users WHERE user_role = 'employee';

    -- Optional: delete admin dummy
    DELETE FROM users WHERE id = '3daaf454-5ded-4f5a-8cad-04f19b70ef72';
END $$;
