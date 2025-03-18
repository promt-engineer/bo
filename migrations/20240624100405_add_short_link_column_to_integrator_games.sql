-- +goose Up
-- +goose StatementBegin
alter table integrator_games add column short_link bool DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table integrator_games drop column short_link;
-- +goose StatementEnd
