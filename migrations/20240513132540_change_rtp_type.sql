-- +goose Up
-- +goose StatementBegin
ALTER TABLE games ALTER COLUMN rtp TYPE INTEGER USING (NULLIF(rtp, '')::NUMERIC::INTEGER);
ALTER TABLE integrator_games ALTER COLUMN rtp TYPE INTEGER USING (NULLIF(rtp, '')::NUMERIC::INTEGER);
UPDATE games SET available_rtp = ARRAY(SELECT NULLIF(value, '')::NUMERIC::INTEGER FROM unnest(available_rtp) AS value);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE games ALTER COLUMN rtp TYPE VARCHAR;
ALTER TABLE integrator_games ALTER COLUMN rtp TYPE VARCHAR;
UPDATE games SET available_rtp = ARRAY(SELECT value FROM unnest(available_rtp) AS value);
-- +goose StatementEnd