-- +goose Up
-- +goose StatementBegin
INSERT INTO lessons (id, name, description, content, "order", chapter_id, created_at, updated_at) VALUES
                                                                                                      (1, 'Hello World', 'Your first Go program.', 'package main\nimport "fmt"\nfunc main() {\n fmt.Println("Hello, world!")\n}', 1, 1, NOW(), NOW()),
                                                                                                      (2, 'If Statements', 'Conditional logic in Go.', 'if condition {\n  // code\n}', 1, 2, NOW(), NOW()),
                                                                                                      (3, 'PostgreSQL Indexes', 'Using indexes to speed up queries.', 'CREATE INDEX idx_name ON table(column);', 1, 3, NOW(), NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM lessons WHERE id IN (1, 2, 3);
-- +goose StatementEnd
