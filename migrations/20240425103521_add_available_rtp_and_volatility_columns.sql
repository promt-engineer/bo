-- +goose Up
-- +goose StatementBegin
alter table games add column available_rtp varchar[];
alter table games add column available_volatility varchar[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table games drop column available_rtp;
alter table games drop column available_volatility;
-- +goose StatementEnd
