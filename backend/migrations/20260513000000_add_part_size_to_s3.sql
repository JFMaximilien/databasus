-- +goose Up
-- +goose StatementBegin
ALTER TABLE s3_storages
    ADD COLUMN s3_part_size BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE s3_storages
    DROP COLUMN s3_part_size;
-- +goose StatementEnd
