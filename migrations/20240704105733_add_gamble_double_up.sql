-- +goose Up
-- +goose StatementBegin
alter table games add column gamble_double_up int default 0 not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table games drop column gamble_double_up;
-- +goose StatementEnd
