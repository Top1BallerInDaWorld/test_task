-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS clicks
(
    timestamp TIMESTAMP        DEFAULT NOW(),
    banner_id INTEGER NOT NULL,
    count     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_clicks_banner_id_timestamp_covering
    ON clicks (banner_id, timestamp) INCLUDE (count);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clicks;
-- +goose StatementEnd
