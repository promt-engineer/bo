-- +goose Up
-- +goose StatementBegin
alter table integrator_games add column rtp varchar;
alter table integrator_games add column volatility varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table integrator_games drop column rtp;
alter table integrator_games drop column volatility;
-- +goose StatementEnd