DO $$
BEGIN
    -- ROLE_ADMIN
    IF NOT EXISTS (
        SELECT 1 FROM keycloak_role
        WHERE name = 'ROLE_ADMIN'
        AND realm_id = (SELECT id FROM realm WHERE name = 'lms')
    ) THEN
        INSERT INTO keycloak_role (
            id, client_realm_constraint, client_role, description,
            name, realm_id, client
        ) VALUES (
            gen_random_uuid(),
            (SELECT id FROM realm WHERE name = 'lms'),
            false,
            'System Administrator',
            'ROLE_ADMIN',
            (SELECT id FROM realm WHERE name = 'lms'),
            NULL
        );
END IF;

    -- ROLE_TEACHER
    IF NOT EXISTS (
        SELECT 1 FROM keycloak_role
        WHERE name = 'ROLE_TEACHER'
        AND realm_id = (SELECT id FROM realm WHERE name = 'lms')
    ) THEN
        INSERT INTO keycloak_role (
            id, client_realm_constraint, client_role, description,
            name, realm_id, client
        ) VALUES (
            gen_random_uuid(),
            (SELECT id FROM realm WHERE name = 'lms'),
            false,
            'Teacher role',
            'ROLE_TEACHER',
            (SELECT id FROM realm WHERE name = 'lms'),
            NULL
        );
END IF;

    -- ROLE_STUDENT
    IF NOT EXISTS (
        SELECT 1 FROM keycloak_role
        WHERE name = 'ROLE_STUDENT'
        AND realm_id = (SELECT id FROM realm WHERE name = 'lms')
    ) THEN
        INSERT INTO keycloak_role (
            id, client_realm_constraint, client_role, description,
            name, realm_id, client
        ) VALUES (
            gen_random_uuid(),
            (SELECT id FROM realm WHERE name = 'lms'),
            false,
            'Student role',
            'ROLE_STUDENT',
            (SELECT id FROM realm WHERE name = 'lms'),
            NULL
        );
END IF;
END $$;

DO $$
DECLARE
v_realm_id TEXT;
    v_admin_id TEXT;
    v_teacher_id TEXT;
    v_student_id TEXT;
BEGIN
    -- Get realm ID dynamically
SELECT id INTO v_realm_id FROM realm WHERE name = 'lms';

-- Create admin user if not exists
IF NOT EXISTS (SELECT 1 FROM user_entity WHERE username = 'admin' AND realm_id = v_realm_id) THEN
        v_admin_id := gen_random_uuid();

INSERT INTO user_entity (
    id, email, email_constraint, email_verified, enabled,
    federation_link, first_name, last_name, realm_id, username,
    created_timestamp, service_account_client_link, not_before
) VALUES (
             v_admin_id,
             'admin@lms.edu', 'admin@lms.edu', true, true,
             NULL, 'System', 'Administrator', v_realm_id, 'admin',
             EXTRACT(EPOCH FROM NOW()) * 1000, NULL, 0
         );

-- For admin user (password: admin123)
-- INSERT INTO credential (
--     id, salt, type, user_id, created_date, user_label,
--     secret_data, credential_data, priority
-- ) VALUES (
--              gen_random_uuid(), '', 'password', v_admin_id, EXTRACT(EPOCH FROM NOW()) * 1000, NULL,
--              '{"value":"$pbkdf2-sha256$i=27500,l=64$RG3gMleD+AvzowLzEEsT5Q$04YQ7nY+G+TzVXq2yhtckptjRt6OhZozXJ5Kq7US8zQ=","salt":""}',
--              '{"hashIterations":27500,"algorithm":"pbkdf2-sha256"}', 10
--          );

-- Assign ROLE_ADMIN
INSERT INTO user_role_mapping (role_id, user_id)
SELECT id, v_admin_id FROM keycloak_role
WHERE name = 'ROLE_ADMIN' AND realm_id = v_realm_id;
END IF;

    -- Create teacher user if not exists
    IF NOT EXISTS (SELECT 1 FROM user_entity WHERE username = 'teacher1' AND realm_id = v_realm_id) THEN
        v_teacher_id := gen_random_uuid();

INSERT INTO user_entity (
    id, email, email_constraint, email_verified, enabled,
    federation_link, first_name, last_name, realm_id, username,
    created_timestamp, service_account_client_link, not_before
) VALUES (
             v_teacher_id,
             'johnsmith@lms.edu', 'johnsmith@lms.edu', true, true,
             NULL, 'John', 'Smith', v_realm_id, 'TeacherJohn',
             EXTRACT(EPOCH FROM NOW()) * 1000, NULL, 0
         );

-- For teacher1 (password: teacher123)
-- INSERT INTO credential (
--     id, salt, type, user_id, created_date, user_label,
--     secret_data, credential_data, priority
-- ) VALUES (
--              gen_random_uuid(), '', 'password', v_teacher_id, EXTRACT(EPOCH FROM NOW()) * 1000, NULL,
--              '{"value":"$2y$12$MuysSNwkv9Kaki8Udavewum4w6xPViHI7oDWB6a.1.seY7bHYQNBS","salt":""}',
--              '{"hashIterations":27500,"algorithm":"pbkdf2-sha256"}', 10
--          );

-- Assign ROLE_TEACHER
INSERT INTO user_role_mapping (role_id, user_id)
SELECT id, v_teacher_id FROM keycloak_role
WHERE name = 'ROLE_TEACHER' AND realm_id = v_realm_id;
END IF;

    -- Create student user if not exists
    IF NOT EXISTS (SELECT 1 FROM user_entity WHERE username = 'student1' AND realm_id = v_realm_id) THEN
        v_student_id := gen_random_uuid();

INSERT INTO user_entity (
    id, email, email_constraint, email_verified, enabled,
    federation_link, first_name, last_name, realm_id, username,
    created_timestamp, service_account_client_link, not_before
) VALUES (
             v_student_id,
             'alicejohnson@lms.edu', 'alicejohnson@lms.edu', true, true,
             NULL, 'Alice', 'Johnson', v_realm_id, 'StudentAlice',
             EXTRACT(EPOCH FROM NOW()) * 1000, NULL, 0
         );

-- Set password (bcrypt hash of 'student123')
-- For student1 (password: student123)
-- INSERT INTO credential (
--     id, salt, type, user_id, created_date, user_label,
--     secret_data, credential_data, priority
-- ) VALUES (
--              gen_random_uuid(), '', 'password', v_student_id, EXTRACT(EPOCH FROM NOW()) * 1000, NULL,
--              '{"value":"$2y$12$cKk.vIYCKShMF6USbk8pIuN4YpL6dW9DIe08rkTlcTI2oOZ6Ek1Cm","salt":""}',
--              '{"hashIterations":27500,"algorithm":"pbkdf2-sha256"}', 10
--          );
-- Assign ROLE_STUDENT
INSERT INTO user_role_mapping (role_id, user_id)
SELECT id, v_student_id FROM keycloak_role
WHERE name = 'ROLE_STUDENT' AND realm_id = v_realm_id;
END IF;
END $$;