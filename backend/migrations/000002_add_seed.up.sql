-- Create user
INSERT INTO users (id, name, email, password)
VALUES (
    uuid_generate_v4(),
    'Test User',
    'test@example.com',
    '$2a$12$6I5s7y7nyXTcK1Axu5Gu4uP9SqKLtOjWIMKMjVLSUH4/bimqe5UPK' -- password: password123
);

-- Create project
INSERT INTO projects (id, name, description, owner_id)
VALUES (
    uuid_generate_v4(),
    'Sample Project',
    'Demo project',
    (SELECT id FROM users LIMIT 1)
);

-- Create tasks
INSERT INTO tasks (id, title, status, priority, project_id)
VALUES
(uuid_generate_v4(), 'Task 1', 'todo', 'low', (SELECT id FROM projects LIMIT 1)),
(uuid_generate_v4(), 'Task 2', 'in_progress', 'medium', (SELECT id FROM projects LIMIT 1)),
(uuid_generate_v4(), 'Task 3', 'done', 'high', (SELECT id FROM projects LIMIT 1));