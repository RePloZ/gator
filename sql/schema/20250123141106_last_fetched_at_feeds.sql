-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds ADD last_fetched_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds DROP last_fetched_at;
-- +goose StatementEnd
