-- +goose Up
-- +goose StatementBegin
alter table integrator_games add column wager_set_id uuid;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table integrator_games drop column wager_set_id;
-- +goose StatementEnd