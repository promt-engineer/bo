-- +goose Up
-- +goose StatementBegin
alter table currency_multipliers add column synonym varchar;
update currency_multipliers set synonym = title;
alter table currency_multipliers alter column synonym SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table currency_multipliers drop column synonym;
-- +goose StatementEnd
