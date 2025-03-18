-- +goose Up
-- +goose StatementBegin
insert into permissions (name, description, subject, endpoint, action)

values ('Get integrator games wager sets', 'Get integrator games wager sets', 'backoffice', '/organizations/:id/wager_set', 'VIEW'),
       ('Update integrator game wager set', 'Update integrator game wager set', 'backoffice', '/organizations/:id/wager_set', 'EDIT'),
       ('Add integrator games wager sets', 'Add integrator games wager sets', 'backoffice', '/organizations/:id/wager_set', 'CREATE'),
       ('Delete integrator games wager sets', 'Delete integrator games wager sets', 'backoffice', '/organizations/:id/wager_set', 'DELETE')  ON CONFLICT (name) DO

UPDATE
    SET
    description = EXCLUDED.description,
    subject = EXCLUDED.subject,
    endpoint = EXCLUDED.endpoint,
    action = EXCLUDED.action;

call refresh_admin_permissions();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
