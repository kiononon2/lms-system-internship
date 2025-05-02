-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO courses (id, name, description, created_at, updated_at) VALUES
                                                                        (1, 'Go Programming Basics', 'An introduction to the Go programming language.', NOW(), NOW()),
                                                                        (2, 'Advanced PostgreSQL', 'Deep dive into PostgreSQL features and performance.', NOW(), NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM courses WHERE id IN (1, 2);

-- +goose StatementEnd
