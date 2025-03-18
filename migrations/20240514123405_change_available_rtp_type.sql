-- +goose Up
-- +goose StatementBegin
ALTER TABLE games ADD COLUMN new_available_rtp INTEGER[];

UPDATE games
SET new_available_rtp = ARRAY(
    SELECT CAST(value AS INTEGER) FROM UNNEST(available_rtp) AS value
);

ALTER TABLE games DROP COLUMN available_rtp;

ALTER TABLE games RENAME COLUMN new_available_rtp TO available_rtp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE games ADD COLUMN new_available_rtp VARCHAR[];
UPDATE games
SET new_available_rtp = ARRAY_TO_STRING(available_rtp, ',');

ALTER TABLE games DROP COLUMN available_rtp;

ALTER TABLE games RENAME COLUMN new_available_rtp TO available_rtp;
-- +goose StatementEnd
