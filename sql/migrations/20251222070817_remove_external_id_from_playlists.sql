-- +goose Up
-- +goose StatementBegin
ALTER TABLE playlists DROP COLUMN external_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE playlists ADD COLUMN external_id TEXT NOT NULL;
-- +goose StatementEnd