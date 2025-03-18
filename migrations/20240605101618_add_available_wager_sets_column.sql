-- +goose Up
-- +goose StatementBegin
ALTER TABLE games ADD COLUMN available_wager_sets_id UUID[];
UPDATE games SET available_wager_sets_id = ARRAY[wager_set_id]::UUID[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE games DROP COLUMN available_wager_sets_id;
-- +goose StatementEnd
