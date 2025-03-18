-- +goose Up
-- +goose StatementBegin
alter table games add column rtp varchar;
alter table games add column volatility varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table games drop column rtp;
alter table games drop column volatility;
-- +goose StatementEnd