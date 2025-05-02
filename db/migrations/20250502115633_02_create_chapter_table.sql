-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO chapters (id, name, description, "order", course_id, created_at, updated_at) VALUES
                                                                                             (1, 'Introduction to Go', 'Getting started with Go.', 1, 1, NOW(), NOW()),
                                                                                             (2, 'Control Structures', 'Learn about if, for, and switch.', 2, 1, NOW(), NOW()),
                                                                                             (3, 'Indexes and Constraints', 'Understanding performance tuning.', 1, 2, NOW(), NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM chapters WHERE id IN (1, 2, 3);
-- +goose StatementEnd
