-- +goose Up
-- +goose StatementBegin
alter table currencies add column base_currency varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table currencies drop column base_currency;
-- +goose StatementEnd