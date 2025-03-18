-- +goose Up
-- +goose StatementBegin
alter table games add column online_volatility bool DEFAULT false;
update games set online_volatility = false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table games drop column online_volatility;
-- +goose StatementEnd
